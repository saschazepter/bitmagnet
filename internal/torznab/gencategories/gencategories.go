package main

import (
  _ "embed"
  "os"
  "path"
  "runtime"
  "strconv"
  "strings"

  "github.com/bitmagnet-io/bitmagnet/internal/maps"
  "github.com/bitmagnet-io/bitmagnet/internal/torznab"
)

// Taken from https://torznab.github.io/spec-1.3-draft/external/newznab/api.html#predefined-categories
//
//go:embed categories.csv
var categoriesCsvString string

func main() {
  categoriesMap, categoriesMapErr := readCategoriesMap()
  checkErr(categoriesMapErr)
  var varNames []struct {
    name string
    id   int
  }
  maxVarNameLength := 0
  var topLevelNames []string
  out := "// Code generated by gencategories. DO NOT EDIT.\n\n"
  out += "package torznab\n\n"
  out += "var categoriesMap = map[int]Category{\n"
  for _, category := range categoriesMap.Values() {
    out += "  " + strconv.Itoa(category.ID) + ": {\n"
    out += "    ID: " + strconv.Itoa(category.ID) + ",\n"
    out += "    Name: \"" + category.Name + "\",\n"
    out += "    Subcat: []Subcategory{\n"
    for _, subcategory := range category.Subcat {
      out += "      {\n"
      out += "        ID: " + strconv.Itoa(subcategory.ID) + ",\n"
      out += "        Name: \"" + subcategory.Name + "\",\n"
      out += "      },\n"
    }
    out += "    },\n"
    out += "  },\n"
    varName := "Category" + strings.Replace(category.Name, "/", "", -1)
    varNames = append(varNames, struct {
      name string
      id   int
    }{name: varName, id: category.ID})
    if len(varName) > maxVarNameLength {
      maxVarNameLength = len(varName)
    }
    if category.ID%1000 == 0 {
      topLevelNames = append(topLevelNames, varName)
    }
  }
  out += "}\n\n"
  out += "var (\n"
  for _, varName := range varNames {
    out += "  " + varName.name
    for i := 0; i < maxVarNameLength-len(varName.name); i++ {
      out += " "
    }
    out += " = categoriesMap[" + strconv.Itoa(varName.id) + "]\n"
  }
  out += ")\n\n"
  out += "var TopLevelCategories = []Category{\n"
  for _, topLevelName := range topLevelNames {
    out += "  " + topLevelName + ",\n"
  }
  out += "}\n"
  _, filename, _, _ := runtime.Caller(0)
  outFile := path.Dir(path.Dir(filename)) + "/categories.gen.go"
  f, fErr := os.Create(outFile)
  checkErr(fErr)
  _, wErr := f.WriteString(out)
  checkErr(wErr)
}

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}

func readCategoriesMap() (maps.InsertMap[int, torznab.Category], error) {
  categoriesMap := maps.NewInsertMap[int, torznab.Category]()
  csvLines := strings.Split(categoriesCsvString, "\n")[1:]
  var currentCategoryID int
  for _, line := range csvLines {
    if len(line) == 0 {
      continue
    }
    parts := strings.Split(line, ",")
    enabled := parts[2][0:1]
    if enabled != "1" {
      continue
    }
    id, idErr := strconv.Atoi(parts[0])
    if idErr != nil {
      return categoriesMap, idErr
    }
    name := parts[1]
    categoriesMap.Set(id, torznab.Category{
      ID:     id,
      Name:   name,
      Subcat: make([]torznab.Subcategory, 0),
    })
    if parts[0][1:] == "000" {
      currentCategoryID = id
    } else {
      currentCategory, _ := categoriesMap.Get(currentCategoryID)
      currentCategory.Subcat = append(currentCategory.Subcat, torznab.Subcategory{
        ID:   id,
        Name: name,
      })
      categoriesMap.Set(currentCategoryID, currentCategory)
    }
  }
  return categoriesMap, nil
}
