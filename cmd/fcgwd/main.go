
package main

import (
    "os"
    "os/signal"
    "syscall"
    "github.com/jcsp/fcgw"
)

func main() {
    st := fcgw.NewSeriesTable()
    sv := fcgw.NewInfluxServer(st)

    st.Logger.Printf("Opening server...")
    sv.Open()

    // HTTP now serving in the background.
    // Wait for a ctrl-c
    signalCh := make(chan os.Signal, 1)
    signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
    select {
    case <-signalCh:
        go func() {
            st.Logger.Printf("Closing server...")
            sv.Close()
        }()
    }
}

