package main

import (
	"fmt"
	"strconv"
	//"io/fs"
	//"flag"
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
func timestamp() string {
   return time.Now().Format("2006-01-02T15:05")
}

func note_fname() string {
   return fmt.Sprintf("%s.md", timestamp())//.Now().Format("2006-01-02T15:05"))
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

func rate(person string, rating string, score int, comment string) error {
   path_to_ppl_folder,err := path_to_folder("")
   if err != nil {
      return err
   }
   filename := fmt.Sprintf("%s/ratings.csv",path_to_ppl_folder)
   f,err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
   defer f.Close()
   if _,err := f.WriteString(fmt.Sprintf("%s,%s,%s,%v,%s\n", timestamp(), person, rating, score, comment)) ; err != nil {
      return err
   }
   return nil
}

func usage(err error){
   if err != nil {
      fmt.Printf("Error: %s", err)
   }
   fmt.Printf("Usage: ppl [write | list] [person name]\n")
   os.Exit(1)
}

func require_narg(number_or_arguments_required int, args []string) error {
   if narg := len(args)-2; narg != number_or_arguments_required {
      return fmt.Errorf("Requires %v arguments, got %v: %s\n", number_or_arguments_required, narg, args)
   }
   return nil
}

func main(){
   //if narg := len(os.Args); narg != 3 {
   //   usage(nil)
   //}
   command_name := os.Args[1]

   switch command_name {
      case "write":
         err := require_narg(1, os.Args)
         if err != nil {
            usage(err)
            return
         }
         err = take_note(os.Args[2])
         if err != nil {
            panic(err)
         }
      case "list":
         err := require_narg(1, os.Args)
         if err != nil {
            usage(err)
            return
         }
         err = list_notes(os.Args[2])
         if err != nil {
            panic(err)
         }
      case "rate":
         //comment := flag.String("m", "", "Comment for this rating")
         //flag.Parse()
         args := os.Args
         err := require_narg(3, args)
         if err != nil {
            usage(err)
            return
         }
         score,err := strconv.Atoi(args[4])
         if err != nil {
            usage(fmt.Errorf("Failed to parse score as integer: %s\n", args[4]))
            return
         }
         err = rate(args[2], args[3], score, "")
         if err != nil {
            panic(err)
         }
      default:
         usage(nil)
   }
}
