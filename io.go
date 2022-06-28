package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "bytes"
    "io/ioutil"
    "time"
)

func (s *Scanner) PrintOutput() {

    writeOutput(Hilite, "\nOUTPUT:\n")
    for i, host := range s.hosts {
        if host.state != Up {
            writeOutput(Hilite, "❌ Host " + strconv.Itoa(i) + ": " + host.ip + " (" + host.name + ")")
            writeOutput(Warn, "  " + host.ip + " is not pingable\n")
            continue;
        }
        httpStr := ""
        writeOutput(Hilite, "✅ Host " + strconv.Itoa(i) + ": " + host.ip + " (" + host.name + ")")
        fmt.Printf("  Open Ports:\n")
        for _, port := range host.ports {
            for _, h := range host.http {
                if h == port { httpStr = httpStr + "(http)" }
            }
            for _, hs := range host.https {
                if hs == port { httpStr = httpStr + "(https)" }
            }
            // TODO make http string colorized
            fmt.Printf("    %d  \t%s    \t%s\n", port, DescribePort(port), httpStr)
        }
        fmt.Println()
        fmt.Println(host.getCmdStrings())
    }
}


func readFile(path string) []string {
    var lines []string

    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
        return lines
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
        return lines
    }
    return lines
}

func (s *Scanner) CreateWorkingDirs() {

    for _, host := range s.hosts {
        dir := host.name
        path := dir + "/nmap"
        if err := os.MkdirAll(path, os.ModePerm); err != nil {
            log.Fatal(err)
        }
        copyTemplate(host)
    }
}

func copyTemplate(host Host) {

    newFile := host.name + "/" + host.name + ".md"
    ogFile := templateFileStr

    t := time.Now()
    dt := t.Format("January 2, 2006")
    ports := host.getPortString()
    
    input, err := ioutil.ReadFile(ogFile)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    cmdStr := host.getCmdStrings()

    output := bytes.Replace(input, []byte("{{date}}"), []byte(dt), -1)
    output = bytes.Replace(output, []byte("{{title}}"), []byte(host.name), -1)
    output = bytes.Replace(output, []byte("{{ip}}"), []byte(host.ip), -1)
    output = bytes.Replace(output, []byte("{{ports}}"), []byte(ports), -1)
    output = bytes.Replace(output, []byte("{{cmds}}"), []byte(cmdStr), -1)

    if err = ioutil.WriteFile(newFile, output, 0666); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

}

func writeOutput(outputType MsgType, text string) {

    switch outputType {
    case Info:
        fmt.Println(text)
    case Hilite:
        fmt.Println(cHilite(text))
    case Warn:
        fmt.Println(cWarn(text))
    case Fatal:
        fmt.Println(cFatal(text))
    default:
        fmt.Println(cFatal("ERROR: writing output"))
        // TODO handle this
    }
}