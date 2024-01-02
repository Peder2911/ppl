package main

import (
	"fmt"
	//"io/fs"
        "io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
        "time"
)

func path_to_folder(subpath string) (string, error) {
   // Returns a folder path within the ppl directory (~/ppl)
   // Also ensures that the folder exists.
   var path_to_folder string
   user,err := user.Current()
   if err != nil {
      return path_to_folder, err
   }
   path_to_folder = path.Join(user.HomeDir, "ppl", subpath)
   os.MkdirAll(path_to_folder, 0700) 
   return path_to_folder, nil
}

func note_fname() string {
   return fmt.Sprintf("%s.md", time.Now().Format("2006-01-02T15:05"))
}

func take_note(person string) error {
   path_to_person_folder,err := path_to_folder(person)
   if err != nil {
      return err
   }
   cmd := exec.Command("nvim", path.Join(path_to_person_folder, note_fname()))
   cmd.Stdin = os.Stdin
   cmd.Stdout = os.Stdout
   cmd.Stderr = os.Stderr
   res := cmd.Run()
   if res != nil {
      fmt.Println(res)
   }
   return nil
}

func list_notes(person string) error {
   path_to_person_folder,err := path_to_folder(person)
   if err != nil {
      return err
   }
   files, err := os.ReadDir(path_to_person_folder)
   for _,file := range files {
      if file.Type().IsRegular() {
         notes_file_path := path.Join(path_to_person_folder, file.Name())
         fmt.Printf("Reading from %s\n", notes_file_path)
         contents,err := ioutil.ReadFile(notes_file_path)
         if err != nil {
            return err
         }
         fmt.Printf("%s\n---\n", string(contents))
      }
   }
   return nil
}

func usage(){
   fmt.Printf("Usage: ppl [write | list] [person name]\n")
   os.Exit(1)
}

func main(){
   if narg := len(os.Args); narg != 3 {
      usage()
   }
   command_name := os.Args[1]

   switch command_name {
      case "write":
         err := take_note(os.Args[2])
         if err != nil {
            panic(err)
         }
      case "list":
         err := list_notes(os.Args[2])
         if err != nil {
            panic(err)
         }
      default:
         usage()
   }
}
