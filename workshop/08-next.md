# 8. Make the factory yours

You now have a project on your laptop, isolated workers, task access, reusable
capabilities and a team that can stop for your answer. Choose the next extension
based on the work you want to delegate.

## Use your own project

Follow [apply the factory to another project](../chapters/07-factory/README.md).
It starts from the configuration you have now: `MODE=mcp`, the ACR and Herdr kits,
a working team and a Beans adapter that can append notes. It shows how to select
a different application path and write a task that asks the team to discover the
project's own setup. Keep its existing guidance rather than replacing it blindly.

## Try another provider when you have a reason

Different roles can use the same assistant and model. If you want to compare
another model's review or implementation, the [provider comparison](../chapters/03-pi/MIXED-MODELS.md)
explains the extra runtime and authentication setup. Read its prerequisites before
changing your working team. Subscription access to one assistant does not imply
API access to another provider.

## Put a dispatcher on the host

So far, every `crew` command ran inside SBX. You could give a single agent a task
with one prompt, then give the team a task through its coordinator. A host agent
is useful when you want to choose among several environments, dispatch a backlog,
or collect results without entering each sandbox yourself.

That is a further integration to build, not a prerequisite for this factory.
Define what the dispatcher may launch, what task and project each worker gets,
and how it returns a commit, check results and questions. Keep execution in SBX;
keep access changes and product decisions visible to you. Start with one task and
one worker before adding concurrent jobs.

Your `factory/sbxenv.yaml`, kits, task requirements and role briefs are the pieces
to reuse. Change one, create a fresh environment, and verify what the next worker
actually receives.

[Back to the workshop map](README.md).
