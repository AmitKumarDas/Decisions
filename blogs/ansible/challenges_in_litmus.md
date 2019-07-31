## Challenges w/ Ansible : Re-assess framework


- Asynchronous task execution in ansible (do things in parallel)
- Ability to obtain on-going stdout of long running tasks, before task completion
- Ability to do diffs of json files/blocks
- Ability to execute "blocks of code" with complex conditionals (ansible can do this against an "imported" taskfile)
- Methods to perform an unit/bdd "test" of ansible playbooks
- Automated quality checks/bots against yaml/playbooks for github repos
- Method to speed up execution of litmusbooks
- Method to achieve idempotency using shell tasks (shell module)
- Methods to Speeden up execution of ansible task 
- Lines of code (yaml) in a playbook (can templating everything help) 


## Mitigation/possible approaches

### Asynchronous task execution in ansible

#### Example use-case No.1:

In the following litmusbook (playbook): https://github.com/litmuschaos/litmus/blob/master/chaoslib/openebs/jiva_controller_network_delay.yaml, 
the task-1 should run asynchronously, allowing task-2 to be executed when task-1 is still running:

- Task-1: Injection of delay on jiva ctrl

```yaml
- name: Inject egress delay of {{network_delay}}ms on jiva controller for {{ chaos_duration }}s
  shell: >
    kubectl exec {{ pumba_pod.stdout }} -n {{ app_ns }} 
    -- pumba netem --interface eth0 --duration {{ chaos_duration }}s delay
    --time {{ network_delay }} re2:k8s_{{ jiva_controller_name[:-1] }} 
  args:
    executable: /bin/bash
```

- Task-2: Verify replicas getting disconnected

```yaml
- name: Verifying the Replica getting disconnected
  shell: >
   kubectl exec -it {{ jiva_controller_pod.stdout }} -n {{ app_ns }} 
   -c {{ jiva_controller_name }} curl http://"{{controller_svc.stdout}}":9501/v1/volumes | jq -r '.data[].replicaCount'
  args:
    executable: /bin/bash
  register: resp
  until: resp.stdout != rcount_before.stdout
  retries: 10
  delay: 15
```

#### Potential Solution:

Use the ansible async in a fire-and-forget mode with poll interval set to 0, and use async_status to check back on completion of this task:

```yaml
- name: Inject egress delay of {{network_delay}}ms on jiva controller for {{ chaos_duration }}s
  shell: >
    kubectl exec {{ pumba_pod.stdout }} -n {{ app_ns }} 
    -- pumba netem --interface eth0 --duration {{ chaos_duration }}s delay
    --time {{ network_delay }} re2:k8s_{{ jiva_controller_name[:-1] }} 
  args:
    executable: /bin/bash
  async: 600 #(chaos_duration + 10)s 
  poll: 0
  register: pumba_chaos

- name: Verifying the Replica getting disconnected
  shell: >
   kubectl exec -it {{ jiva_controller_pod.stdout }} -n {{ app_ns }} 
   -c {{ jiva_controller_name }} curl http://"{{controller_svc.stdout}}":9501/v1/volumes | jq -r '.data[].replicaCount'
  args:
    executable: /bin/bash
  register: resp
  until: resp.stdout != rcount_before.stdout
  retries: 10
  delay: 15

- name: 'check on pumba chaos task'
  async_status:
    jid: "{{ pumba_chaos.ansible_job_id }}"
  register: chaos_result
  until: chaos_result.finished
  retries: 30
```

### Automated quality checks/bots against yaml/playbooks for github repos

- Building your own bot to implement/enforce checks into your playbooks
 
  - http://willthames.github.io/2016/06/28/announcing-ansible-review.html
  - https://funinit.wordpress.com/2018/12/19/how-to-develop-ansible-review-standards/

- What is the nature of static checks ? Are they helpful in context of e2e code? 

###  Methods to perform an unit/bdd "test" of ansible playbooks (functional test) 

- Use molecule to test functional requirement of every util
  - https://www.jeffgeerling.com/blog/2018/testing-your-ansible-roles-molecule

- Set-up CI/CD (say, bring up cluster with predefined dependencies for openebs) 
  - Execute litmusbook
