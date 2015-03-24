# Command line arguments for the flue-agent

I think I should implement the following command line arguments

* `flue -action=status -uuid=uuid`
* `flue -acion=add -task=layerX.tasks [-uuid=uuid`]
This action will add all the tasks from layerX.tasks to the *stack*.
if **uuid** is specified, connect to the running **flue** instance and add the tasks to its queue.
If **uuid** is not specified, create a new server instance and add the command line arguments
* `flue -action=list` to list all the running uuid


