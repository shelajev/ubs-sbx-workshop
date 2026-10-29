import importlib.machinery
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch


CREW = Path(__file__).resolve().parents[2] / 'chapters/support/bin/crew'


class CrewNotificationTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory(prefix='workshop-crew-test-')
        self.addCleanup(self.directory.cleanup)
        environment = patch.dict(os.environ, {'FACTORY_DIR': self.directory.name})
        environment.start()
        self.addCleanup(environment.stop)
        loader = importlib.machinery.SourceFileLoader('crew_under_test', str(CREW))
        spec = importlib.util.spec_from_loader(loader.name, loader)
        self.crew = importlib.util.module_from_spec(spec)
        loader.exec_module(self.crew)
        (self.crew.base / 'state.json').write_text(json.dumps({'task_id': 'sbx-test'}))
        self.messages = ['msg-old-unread']
        self.claims = {'msg-old-unread'}
        self.notifications = []
        self.prompts = []
        self.busy = False
        runner = patch.object(self.crew.subprocess, 'run', side_effect=self.run_command)
        runner.start()
        self.addCleanup(runner.stop)
        sleeper = patch.object(self.crew.time, 'sleep')
        sleeper.start()
        self.addCleanup(sleeper.stop)

    def run_command(self, command, **kwargs):
        tool = Path(command[0]).name
        if tool == 'handoff':
            self.assertEqual(command[1], 'send')
            self.messages.append('msg-new-saved')
            return subprocess.CompletedProcess(command, 0, stdout='msg-new-saved\n')
        self.assertEqual(tool, 'crew-notify')
        self.notifications.append(command[:])
        message_id = (command[command.index('--message-id') + 1]
                      if '--message-id' in command else self.messages[0])
        if message_id in self.claims:
            return subprocess.CompletedProcess(command, 9, stderr='already notified')
        if self.busy:
            self.busy = False
            return subprocess.CompletedProcess(command, 3, stderr='recipient busy')
        self.claims.add(message_id)
        self.prompts.append(message_id)
        return subprocess.CompletedProcess(command, 0, stderr='')

    def test_new_message_is_prompted_despite_claimed_oldest_message(self):
        self.crew.send('note', 'Review the new commit', target='qa', sender='developer')
        self.assertEqual(self.messages, ['msg-old-unread', 'msg-new-saved'])
        self.assertEqual(self.prompts, ['msg-new-saved'])
        self.assertEqual(self.notifications[0][1:], ['qa', '--message-id', 'msg-new-saved'])

    def test_busy_retry_keeps_the_new_message_identity(self):
        self.busy = True
        self.crew.send('note', 'Review the new commit', target='qa', sender='developer')
        self.assertEqual(self.prompts, ['msg-new-saved'])
        self.assertEqual(len(self.notifications), 2)
        self.assertEqual(self.notifications[0], self.notifications[1])
        self.assertEqual(self.notifications[1][1:], ['qa', '--message-id', 'msg-new-saved'])
        self.assertEqual(self.messages.count('msg-new-saved'), 1)

    def test_submit_notifies_about_the_saved_assignment(self):
        with patch.object(self.crew, 'ready'), patch.object(self.crew.sys, 'argv', ['crew', 'submit']):
            self.crew.main()
        self.assertTrue((self.crew.base / 'submitted').exists())
        self.assertEqual(self.prompts, ['msg-new-saved'])
        self.assertEqual(self.notifications[0][1:], ['coordinator', '--message-id', 'msg-new-saved'])

    def test_deliver_retains_existing_inbox_selection(self):
        self.claims.clear()
        with patch.object(self.crew.sys, 'argv', ['crew', 'deliver']):
            self.crew.main()
        self.assertEqual(self.prompts, ['msg-old-unread'])
        self.assertEqual(self.notifications[0][1:], ['coordinator'])


if __name__ == '__main__':
    unittest.main()
