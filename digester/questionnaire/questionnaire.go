package questionnaire

import (
	"bufio"
	"fmt"

	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/digester/common"

	//"github.com/Microsoft/kunlun/digester/detector"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/Microsoft/kunlun/digester/vmgroupcalc"
	//"strings"
)

var (
	scanner = bufio.NewScanner(os.Stdin)
)

// TODO need a function to construct Q&A

func Run(state storage.State, filePath string) common.Blueprint {
	var err error
	var doFlag bool
	bp, _ := common.ImportBlueprintYaml(filePath)

	fmt.Printf("Project path?")
	if bp.NonInfra.ProjectSourceCodePath != "" {
		fmt.Printf(" Default: %s.", bp.NonInfra.ProjectSourceCodePath)
	}
	fmt.Printf("\n")
	scanner.Scan()
	path := scanner.Text()
	if path != "" {
		bp.NonInfra.ProjectSourceCodePath = path
	}
	/*
		d, err := detector.New(path)
		if err != nil {
			log.Fatal(err)
		}


		possiblePackageManagers := d.DetectPackageManager()
		fmt.Printf("What is the package manager for the project?")
		for i, pm := range possiblePackageManagers {
			if i == 0 {
				fmt.Printf(" %s \\", strings.ToUpper(string(pm)))
			} else {
				fmt.Printf(" %s \\", string(pm))
			}
		}
		if len(possiblePackageManagers) > 0 {
			fmt.Printf(" other?\n")
		} else {
			fmt.Printf("\n")
		}
		scanner.Scan()
		inputPackageManager := scanner.Text()
		if inputPackageManager == "" {
			inputPackageManager = string(possiblePackageManagers[0])
		}
		d.ConfirmPackageManager(inputPackageManager)

		possibleFrameworks := d.DetectFramework()
		fmt.Printf("What is the framework of the project?")
		for i, fw := range possibleFrameworks {
			if i == 0 {
				fmt.Printf(" %s \\", strings.ToUpper(string(fw)))
			} else {
				fmt.Printf(" %s \\", string(fw))
			}
		}
		if len(possibleFrameworks) > 0 {
			fmt.Printf(" other?\n")
		} else {
			fmt.Printf(" NONE?\n")
		}
		scanner.Scan()
		inputFramework := scanner.Text()
		if inputFramework == "" {
			if len(possibleFrameworks) > 0 {
				inputFramework = string(possibleFrameworks[0])
			}
		}
		d.ConfirmFramework(inputFramework)

		d.DetectConfig()

		bp = d.ExposeKnownInfo()
	*/

	// Ask for the empty fields
	bpNonInfra := &bp.NonInfra
	doFlag = true
	for doFlag {
		doFlag = false
		fmt.Printf("What's the programming language?")
		if bpNonInfra.ProgrammingLanguage != "" {
			fmt.Printf(" Default: %s.", bpNonInfra.ProgrammingLanguage)
		}
		fmt.Printf(" Allowed values: {php}.\n")
		scanner.Scan()
		input := scanner.Text()
		if input != "" {
			pl, err := common.ParseProgrammingLanguage(input)
			if err != nil {
				fmt.Println(err)
				doFlag = true
				continue
			}
			bpNonInfra.ProgrammingLanguage = pl
		}
	}

	if len(bpNonInfra.Databases) > 0 {
		needExtraInfo := false
		fmt.Println("Here is the database(s):")
		for i, db := range bpNonInfra.Databases {
			fmt.Printf("No.%d: {\n", i+1)
			s := reflect.ValueOf(&db).Elem()
			for j := 0; j < s.NumField(); j++ {
				valField := s.Field(j)
				typeField := s.Type().Field(j)
				tag := typeField.Tag
				val := valField.Interface()
				if valField.Kind() == reflect.Int {
					fmt.Printf("  %s: %d\n", tag.Get("name"), val)
				} else {
					fmt.Printf("  %s: %s\n", tag.Get("name"), val)
				}
				if val == reflect.Zero(valField.Type()).Interface() {
					needExtraInfo = true
				}
			}
			fmt.Println("}")
		}
		if needExtraInfo {
			fmt.Println("Please help fill the blank field(s).")
			for i := range bpNonInfra.Databases {
				db := &bpNonInfra.Databases[i]
				askForDbEmptyFields(i+1, db)
			}
			fmt.Println("Done.")
		}
	}

	extraDatabasesNum := 0
	doFlag = true
	for doFlag {
		doFlag = false
		if len(bpNonInfra.Databases) > 0 {
			fmt.Println("Do you use any more databases? How many? Default: 0. Allowed values: {n | n >= 0}.")
		} else {
			fmt.Println("Do you use any databases? How many? Default: 0. Allowed values: {n | n >= 0}.")
		}
		scanner.Scan()
		input := scanner.Text()
		if input != "" {
			extraDatabasesNum, err = strconv.Atoi(input)
			if err != nil {
				fmt.Println(err)
				doFlag = true
				continue
			}
		}
	}
	for i := 1; i <= extraDatabasesNum; i++ {
		newDb := common.Database{}
		askForDbEmptyFields(len(bpNonInfra.Databases)+1, &newDb)
		bpNonInfra.Databases = append(bpNonInfra.Databases, newDb)
	}

	// Ask for Misc
	if bp.Misc.ResourceGroupName == "" {
		bp.Misc.ResourceGroupName = state.EnvID
	}
	s := reflect.ValueOf(&bp.Misc).Elem()
	for i := 0; i < s.NumField(); i++ {
		valField := s.Field(i)
		typeField := s.Type().Field(i)
		tag := typeField.Tag
		if valField.Kind() == reflect.Int {
			val := valField.Interface().(int)
			var defaultVal int
			if val == 0 {
				defaultVal, err = strconv.Atoi(tag.Get("default"))
				if err != nil {
					log.Fatal(err)
				}
			} else {
				defaultVal = val
			}
			doFlag = true
			for doFlag {
				doFlag = false
				fmt.Printf(
					"%s Default: %d.\n",
					tag.Get("question"),
					defaultVal,
				)
				scanner.Scan()
				input := scanner.Text()
				if input == "" {
					valField.Set(reflect.ValueOf(defaultVal))
				} else {
					inputToInt, err := strconv.Atoi(input)
					if err != nil {
						fmt.Println(err)
						doFlag = true
						continue
					}
					valField.Set(reflect.ValueOf(inputToInt))
				}
			}
		} else {
			val := valField.Interface().(string)
			var defaultVal string
			if val == "" {
				defaultVal = tag.Get("default")
			} else {
				defaultVal = val
			}
			fmt.Printf(
				"%s Default: %s.\n",
				tag.Get("question"),
				defaultVal,
			)
			scanner.Scan()
			input := scanner.Text()

			if input == "" {
				valField.SetString(defaultVal)
			} else {
				valField.SetString(input)
			}
		}
	}

	bp.Infra = vmgroupcalc.Calc(vmgroupcalc.Requirment{
		ConcurrentUserNumber: bp.Misc.ConcurrentUserNumber,
	})

	return bp
}

func askForDbEmptyFields(num int, db *common.Database) {
	s := reflect.ValueOf(db).Elem()
	for i := 0; i < s.NumField(); i++ {
		valField := s.Field(i)
		typeField := s.Type().Field(i)
		tag := typeField.Tag
		val := valField.Interface()
		allowedSentence := ""
		if tag.Get("allow") != "" {
			allowedSentence = fmt.Sprintf(" Allowed values: %s.", tag.Get("allow"))
		}
		if val == reflect.Zero(valField.Type()).Interface() {
			doFlag := true
			for doFlag {
				doFlag = false
				fmt.Printf(
					"For the database No.%d: %s%s\n",
					num,
					tag.Get("question"),
					allowedSentence,
				)
				scanner.Scan()
				input := scanner.Text()
				if err := db.ValidateField(typeField.Name, input, &valField); err != nil {
					fmt.Println(err)
					doFlag = true
					continue
				}
			}
		}
	}
}
