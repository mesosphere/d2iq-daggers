// Package help provides a help command for mage targets defines in mesosphere/daggers repository.
//
// The purpose of the dedicated package is to avoid name conflicts between different mage targets since namespaces
// is optional during mage import.
//
// The only limitation of this package is that it can only be imported once in a magefile and it should be under the
// help namespace (eg `//mage:import help`).
package help
