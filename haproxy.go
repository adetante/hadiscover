package main

import(
    "text/template"
    "os"
    "os/exec"
    "fmt"
    "strconv"
)

var tpl *template.Template = nil
var pid int = -1

func createConfigFile(backends []Backend, templateFile, outputFile string)(error){
    cfgFile,_ := os.Create(outputFile)
    defer func() {
        cfgFile.Close();
    }()

    if(tpl == nil){
        var err error = nil
        tpl, err = template.ParseFiles(templateFile)
        if (err != nil){
            return err
        }
    }
    
    return tpl.Execute(cfgFile, backends)
}

func reloadHAproxy(command, configFile string)(error){
    var cmd *exec.Cmd = nil
    if pid == -1{
        cmd = exec.Command(command,"-f",configFile)
    } else{
        cmd = exec.Command(command,"-f",configFile,"-sf",strconv.Itoa(pid))
    }
    
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if (err == nil){
        pid = cmd.Process.Pid
        fmt.Println("HAproxy started with pid ",pid)
    }
    return err
}