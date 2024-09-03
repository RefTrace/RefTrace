"""
process FOO {
    cpus 100
}
"""

def rule_check_cpu_directive(module):
    for process in module.processes:
        cpu_directives = process.directives.cpus
        if not cpu_directives:
            # No CPU directive, so we skip this process
            return
        
        for cpu_directive in cpu_directives:
            cpu_value = cpu_directive.num
            if cpu_value < 2 or cpu_value > 96:
                fatal("Process %s has an invalid CPU value. It should be >= 2 and <= 96, but it is %d" % (process.name, cpu_value))