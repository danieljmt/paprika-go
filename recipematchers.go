package paprika

import "strings"

type RecipeMatcher func(*Recipe) bool

func NameContains(name string) RecipeMatcher {
	return func(r *Recipe) bool {
		return strings.Contains(r.Name, name)
	}
}

func HasIngredient(ingredient string) RecipeMatcher {
	return func(r *Recipe) bool {
		return strings.Contains(r.Ingredients, ingredient)
	}
}

func HasCategory(category string) RecipeMatcher {
	return func(r *Recipe) bool {
		for _, c := range r.Categories {
			if c == category {
				return true
			}
		}
		return false
	}
}
