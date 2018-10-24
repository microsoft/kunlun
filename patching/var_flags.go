package apis

type VarFlags struct {
	VarsFiles   []VarsFileArg `long:"vars-file"  short:"l" value-name:"PATH"      description:"Load variables from a YAML file"`
	VarsFSStore VarsFSStore   `long:"vars-store"           value-name:"PATH"      description:"Load/save variables from/to a YAML file"`
}

func (f VarFlags) AsVariables() Variables {
	var firstToUse []Variables

	staticVars := StaticVariables{}

	// for i, _ := range f.VarsEnvs {
	// 	for k, v := range f.VarsEnvs[i].Vars {
	// 		staticVars[k] = v
	// 	}
	// }

	for i, _ := range f.VarsFiles {
		for k, v := range f.VarsFiles[i].Vars {
			staticVars[k] = v
		}
	}

	// for i, _ := range f.VarFiles {
	// 	for k, v := range f.VarFiles[i].Vars {
	// 		staticVars[k] = v
	// 	}
	// }

	// for _, kv := range f.VarKVs {
	// 	staticVars[kv.Name] = kv.Value
	// }

	firstToUse = append(firstToUse, staticVars)

	store := &f.VarsFSStore

	if f.VarsFSStore.IsSet() {
		firstToUse = append(firstToUse, store)
	}

	vars := NewMultiVars(firstToUse)

	if f.VarsFSStore.IsSet() {
		store.ValueGeneratorFactory = NewValueGeneratorConcrete(NewVarsCertLoader(vars))
	}

	return vars
}
