# Command line arguments of the flue-agent

I think I should implement the following command line arguments

* `flue -action=status -uuid=uuid`

Get the running status of the agent identified by **uuid**. The message is return in JSON.

* `flue -acion=add -task=layerX.tasks [-uuid=uuid`]

This action will add all the tasks from layerX.tasks to the *stack*.
if **uuid** is specified, connect to the running **flue** instance and add the tasks to its queue.
If **uuid** is not specified, create a new server instance and add the command line arguments

* `flue -action=list` to list all the running uuid

* `flue -action=wipe -uuid=uuid` to wipe the socket identified by uuid


