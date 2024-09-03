# try with example_process2.nf

def rule_check_input_output_types(module):
    for process in module.processes:
        # Check if the process has exactly one input path
        if len(process.inputs.paths) != 1:
            error("Process '%s' should have exactly one input path" % process.name)
        else:
            print("Process '%s' has the correct number of input paths" % process.name)

        # Check if the process has exactly two output vals
        if len(process.outputs.vals) != 2:
            error("Process '%s' should have exactly two output vals" % process.name)
        else:
            print("Process '%s' has the correct number of output vals" % process.name)

        # Check if the first output val has the correct properties
        if process.outputs.vals and len(process.outputs.vals) > 0:
            first_val = process.outputs.vals[0]
            if first_val.emit != "hello" or not first_val.optional or first_val.topic != "report":
                error("Process '%s' first output val has incorrect properties" % process.name)
            else:
                print("Process '%s' first output val has correct properties" % process.name)

        # If all conditions are met, print a success message
        if (len(process.inputs.paths) == 1 and 
            len(process.outputs.vals) == 2 and 
            process.outputs.vals[0].emit == "hello" and 
            process.outputs.vals[0].optional and 
            process.outputs.vals[0].topic == "report"):
            print("Process '%s' passes all input and output checks" % process.name)
