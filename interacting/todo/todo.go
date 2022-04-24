package todo

import(
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "time"
)

// item struct represents a ToDo item
type item struct {
  Task string
  Done bool
  CreatedAt time.Time
  CompletedAt time.Time
}

// List represents a list of ToDo items
type List []item

// Add creates a new todo item and appends it to the list
func(l *List) Add(task string) {
  t := item{
    Task: task,
    Done: false,
    CreatedAt: time.Now(),
    CompletedAt: time.Time{},
  }

  *l = append(*l,t)

}

// Delete method deletes a ToDo item from the list

func (l *List) Delete(i int) error {
  ls := *l
  if i <= 0 || i > len(ls) {
    return fmt.Errorf("item %d does not exist", i)
  }

  //Adjusting index for 0 based index
  *l = append(ls[:i-1],ls[i:]...)

  return nil

}

// Complete method marks a ToDo item as completed
// by setting Done = true and CompletedAt to the current time

func (l *List) Complete(i int) error {
  ls := *l
  if i<= 0 || i > len(ls) {
    return fmt.Errorf("item %d does not exist",i)
  }

  //Adjusting for 0 based index
  ls[i-1].Done = true
  ls[i-1].CompletedAt = time.Now()

  return nil
}


// Save method encodes the list as JSON and saves it
// using the provided filename
func (l *List) Save(filename string) error {
  js, err := json.Marshal(l)
  if err != nil {
    return err
  }

  return ioutil.WriteFile(filename, js, 0644)
}

// Get method opens the provided file name, decodes
// the JSON data, and parses it into a List
func (l *List) Get(filename string) error{

  file, err := ioutil.ReadFile(filename)
  if err != nil {
    if errors.Is(err, os.ErrNotExist) {
      return nil
    }

    return err
  }

  if len(file) == 0 {
    return nil
  }

  return json.Unmarshal(file, l)
}

// String prints out a formatted list
// Implements the fmt.Stringer interface
func(l *List) String() string {
  formatted := ""

  for k, t := range *l {
    prefix := " "
    if t.Done {
      prefix = "X"
    }
    
    // Adjust the item number k to print numbers starting from 1 instead of 0
    formatted += fmt.Sprintf("%s%d: %s\n",prefix,k+1,t.Task)
  }

  return formatted
}