package pchecker

/**
 * @author  papajuan
 * @date    10/10/2024
 **/

var DefaultFalseNegatives = map[string]bool{
	"asshole": true,
	"dumbass": true, // ass -> bASS (FP) -> dumBASS
	"nigger":  true,
}
