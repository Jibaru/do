package resolver

import "github.com/jibaru/do/internal/types"

type LetResolver interface {
	// Resolve resolves all variables and functions in the let section.
	Resolve(sentences *types.Sentences) (*types.Sentences, error)
}

type letResolver struct {
}

func NewLetResolver() LetResolver {
	return &letResolver{}
}

func (r *letResolver) Resolve(sentences *types.Sentences) (*types.Sentences, error) {
	if sentences == nil {
		return nil, nil
	}

	resolvedSentences := types.NewSentences()

	for _, sentence := range sentences.Entries() {
		key := sentence.Key
		value := sentence.Value

		switch val := value.(type) {
		case types.ReferenceToVariable:
			realValue, ok := resolvedSentences.Get(val.Value)
			if !ok {
				return nil, NewReferenceToVariableNotFoundError(key, val.Value)
			}

			resolvedSentences.Set(key, realValue)
			continue
		case types.Func:
			fn := value.(types.Func)
			if !fn.HasReferences() {
				err := r.ResolveFunction(fn, key, resolvedSentences)
				if err != nil {
					return nil, err
				}
			} else {
				for i, arg := range fn.Args {
					switch arg.(type) {
					case types.ReferenceToVariable:
						realValue, ok := resolvedSentences.Get(arg.(types.ReferenceToVariable).Value)
						if !ok {
							return nil, NewReferenceToVariableNotFoundError(key, arg.(types.ReferenceToVariable).Value)
						}

						fn.Args[i] = realValue
					}
				}

				err := r.ResolveFunction(fn, key, resolvedSentences)
				if err != nil {
					return nil, err
				}
			}
		case types.Map:
			return nil, NewInvalidVariablesError("map type is not allowed in sentences")
		default:
			resolvedSentences.Set(key, value)
		}
	}

	return resolvedSentences, nil
}

func (r *letResolver) ResolveFunction(fn types.Func, key string, resolvedVariables *types.Sentences) error {
	resolvedFn, err := fn.Resolve()
	if err != nil {
		return err
	}

	switch resolvedFn.(type) {
	case types.EnvFunc:
		envFunc := resolvedFn.(types.EnvFunc)
		resolvedVariables.Set(key, envFunc.Exec())
	case types.FileFunc:
		fileFunc := resolvedFn.(types.FileFunc)
		resolvedVariables.Set(key, fileFunc.Exec())
	}

	return nil
}
