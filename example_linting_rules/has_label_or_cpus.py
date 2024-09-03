# This file should exist in the root of your pipeline directory
def has_label(directives):
    return len(directives.label) > 0

def has_cpus(directives):
    return len(directives.cpus) > 0

def rule_has_label_or_cpus(module):
    for process in module.processes:
        if not (has_label(process.directives) or has_cpus(process.directives)):
            fatal("process %s has no label or cpus directive" % process.name)