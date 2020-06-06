package match

default match=false
default exactMatch=false
match{
    prefix :=trim_suffix(input.upperOrderArea,"0")
    startswith(input.workerArea,prefix)
    trace(sprintf("Hello There! %v", [prefix]))
}

exactMatch{
    input.upperOrderArea==input.workerArea
}