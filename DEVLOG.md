# Workflow Automation Server Development Log

<details>
<summary><h3>22/08/2025</h3></summary>

**Work done**
- Limit the task retry delay duration to only have positive values
- Limit number of tasks in a workflow
- Add pointers to domain objects

**TODO**
- Implement HTTP and LOG type tasks payload
- Implement task execution
</details>

<details>
<summary><h3>25/08/2025</h3></summary>

**Work done**
- Add payload types to task in gRPC and domain

**TODO**
- Mask HTTP payload body a map
- Add authentication to HTTP payload
- Implement task execution
</details>

<details>
<summary><h3>08/09/2025</h3></summary>

**Work done**
- Implemented HTTP payload

**TODO**
- Implement task execution
</details>

<details>
<summary><h3>12/09/2025</h3></summary>

**Work done**
- Bug fixes and tested workflow, tasks and payload validations

**TODO**
- Implement task execution
- Implement better error handling
</details>

<details>
<summary><h3>19/09/2025</h3></summary>

**Work done**
- Started implementing workflow execution for log tasks
- Started adding sqlite repository

**TODO**
- Finish workflow log task executing
- Finish implementing the sqlite repo
- Implement pooling for the db
</details>