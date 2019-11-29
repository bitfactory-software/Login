package main

import (
    "io"
    "net/http"
    "os"
    "archive/zip"
    "path/filepath"
    "strings"
    "fmt"
)

func main() {

    exe := "AlexLoginServer/exe"
    zip := "AlexLoginServer/ALEX.zip"
    fileUrl := "https://s3.amazonaws.com/releases.bitfactory.at.development/ALEX64.zip"
    loginServer := "Bfx.Alex.Login.UI.Web.Server.exe"


    fmt.Println( "building directory structure" )
    os.MkdirAll( exe, os.ModePerm )
    
    fmt.Println( "downloading file" )
    if err := DownloadFile(zip, fileUrl); err != nil {
        panic(err)
    }

    fmt.Println( "unzipping" );
    if err := Unzip(zip, exe); err != nil {
        panic(err)
    }

    fmt.Println( "done!" )
    fmt.Println( "start " + exe + "/" + loginServer + " with elevated rights" )
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}

func Unzip(src string, dest string) (error) {

    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer r.Close()

    for _, f := range r.File {

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)

        // Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return fmt.Errorf("%s: illegal file path", fpath)
        }

        if f.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return err
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return err
        }

        rc, err := f.Open()
        if err != nil {
            return err
        }

        _, err = io.Copy(outFile, rc)

        // Close the file without defer to close before next iteration of loop
        outFile.Close()
        rc.Close()

        if err != nil {
            return err
        }
    }
    return nil
}