package pchecker

import (
	"strings"
	"testing"
)

/**
 * @author  papajuan
 * @date    10/9/2024
 **/

func TestProfanityDetector_Censor(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          "what the fuuuuuck",
			expectedOutput: "what the ***",
		},
		{
			input:          "what the N I G   G E R my",
			expectedOutput: "what the ***",
		},
		{
			input:          "what the poop",
			expectedOutput: "what the ***",
		},
		{
			input:          "fuck this",
			expectedOutput: "*** this",
		},
		{
			input:          "one penis, two vaginas, three dicks, four sluts, five whores and a flower",
			expectedOutput: "one ***, two ***, three ***, four ***, five *** and a flower",
		},
		{
			input:          "Censor doesn't support sanitizing '()' into 'o', because it's two characters. Proof: c()ck. Maybe one day I'll have time to fix it.",
			expectedOutput: "Censor doesn't support sanitizing '()' into 'o', because it's two characters. Proof: c()ck. Maybe one day I'll have time to fix it.",
		},
		{
			input:          "fuck shit fuck",
			expectedOutput: "*** *** ***",
		},
		{
			input:          "fuckfuck",
			expectedOutput: "***",
		},
		{
			input:          "fuck this shit",
			expectedOutput: "*** this ***",
		},
		{
			input:          "F   u   C  k th1$ $h!t",
			expectedOutput: "*** th1$ ***",
		},
		{
			input:          "@$$h073",
			expectedOutput: "***",
		},
		{
			input:          "hello, world!",
			expectedOutput: "hello, world!",
		},
		{
			input:          "Hey asshole, are y()u an assassin? If not, fuck off.",
			expectedOutput: "Hey ***, are y()u an assassin? If not, *** off.",
		},
		{
			input:          "I am from Scunthorpe, north Lincolnshire",
			expectedOutput: "I am from Scunthorpe, north Lincolnshire",
		},
		{
			input:          "He is an associate of mine",
			expectedOutput: "He is an associate of mine",
		},
		{
			input:          "But the table is on fucking fire",
			expectedOutput: "But the table is on *** fire",
		},
		{
			input:          "““““““““““““But the table is on fucking fire“",
			expectedOutput: "““““““““““““But the table is on *** fire“",
		},
		{
			input:          "f.u_ck this s.h-i~t",
			expectedOutput: "*** this ***",
		},
		{
			input:          "glass",
			expectedOutput: "glass",
		},
		{
			input:          "ы",
			expectedOutput: "ы",
		},
		{
			input:          "documentdocument", // false positives (https://github.com/TwiN/go-away/issues/30)
			expectedOutput: "documentdocument",
		},
		{
			input:          "dumbassdumbass fuckfuckfuck", // false negatives (https://github.com/TwiN/go-away/issues/30)
			expectedOutput: "*** ***",
		},
		{
			input:          "document fuck document fuck",
			expectedOutput: "document *** document ***",
		},
		{
			input:          "fučk ÄšŚ pÓöp pÉnìŚ bitčh f    u     c k", // accents
			expectedOutput: "*** *** *** *** *** ***",
		},
	}
	for _, tt := range tests {
		t.Run("default_"+tt.input, func(t *testing.T) {
			censored := Censor(tt.input)
			if censored != tt.expectedOutput {
				t.Errorf("expected '%s', got '%s'", tt.expectedOutput, censored)
			}
		})
	}
}

func TestFalsePositives(t *testing.T) {
	sentences := []string{
		"I am from Scunthorpe, north Lincolnshire",
		"He is an associate of mine",
		"Are you an assassin?",
		"But the table is on fire",
		"glass",
		"grass",
		"classic",
		"classification",
		"passion",
		"carcass",
		"cassandra",
		"just push it down the ledge", // puSH IT
		"has steph",                   // hAS Steph
		"was steph",                   // wAS Steph
		"hot water",                   // hoT WATer
		"Phoenix",                     // pHOEnix
		"systems exist",               // systemS EXist
		"saturday",                    // saTURDay
		"therapeutic",
		"press the button",
	}
	t.Run("Test False Positives", func(t *testing.T) {
		for _, s := range sentences {
			if strings.ContainsRune(Censor(s), '*') {
				t.Error("Expected false, got true from:", s)
			}
		}
	})
}

func TestSentencesWithFalsePositivesAndProfanities(t *testing.T) {
	t.Run("Test Sentences With False Positives And Profanities", func(t *testing.T) {
		if s := Censor("You are a associate"); strings.ContainsRune(s, '*') {
			t.Error("Expected true, got false from sentence")
		}
		if s := Censor("Go away, asshole!"); !strings.ContainsRune(s, '*') {
			t.Error("Expected true, got false from sentence", s)
		}
	})
}

// "The Adventures of Sherlock Holmes" by Arthur Conan Doyle is in the public domain,
// which makes it a perfect source to use as reference.
func TestSentencesFromTheAdventuresOfSherlockHolmes(t *testing.T) {
	defaultProfanityDetector = nil
	sentences := []string{
		"I had called upon my friend, Mr. Sherlock Holmes, one day in the autumn of last year and found him in deep conversation with a very stout, florid-faced, elderly gentleman with fiery red hair.",
		"With an apology for my intrusion, I was about to withdraw when Holmes pulled me abruptly into the room and closed the door behind me.",
		"You could not possibly have come at a better time, my dear Watson, he said cordially",
		"I was afraid that you were engaged.",
		"So I am. Very much so.",
		"Then I can wait in the next room.",
		"Not at all. This gentleman, Mr. Wilson, has been my partner and helper in many of my most successful cases, and I have no doubt that he will be of the utmost use to me in yours also.",
		"The stout gentleman half rose from his chair and gave a bob of greeting, with a quick little questioning glance from his small fat-encircled eyes",
		"Try the settee, said Holmes, relapsing into his armchair and putting his fingertips together, as was his custom when in judicial moods.",
		"I know, my dear Watson, that you share my love of all that is bizarre and outside the conventions and humdrum routine of everyday life.",
		"You have shown your relish for it by the enthusiasm which has prompted you to chronicle, and, if you will excuse my saying so, somewhat to embellish so many of my own little adventures.",
		"You did, Doctor, but none the less you must come round to my view, for otherwise I shall keep on piling fact upon fact on you until your reason breaks down under them and acknowledges me to be right.",
		"Now, Mr. Jabez Wilson here has been good enough to call upon me this morning, and to begin a narrative which promises to be one of the most singular which I have listened to for some time.",
		"You have heard me remark that the strangest and most unique things are very often connected not with the larger but with the smaller crimes, and occasionally",
		"indeed, where there is room for doubt whether any positive crime has been committed.",
		"As far as I have heard it is impossible for me to say whether the present case is an instance of crime or not, but the course of events is certainly among the most singular that I have ever listened to.",
		"Perhaps, Mr. Wilson, you would have the great kindness to recommence your narrative.",
		"I ask you not merely because my friend Dr. Watson has not heard the opening part but also because the peculiar nature of the story makes me anxious to have every possible detail from your lips.",
		"As a rule, when I have heard some slight indication of the course of events, I am able to guide myself by the thousands of other similar cases which occur to my memory.",
		"In the present instance I am forced to admit that the facts are, to the best of my belief, unique.",
		"We had reached the same crowded thoroughfare in which we had found ourselves in the morning.",
		"Our cabs were dismissed, and, following the guidance of Mr. Merryweather, we passed down a narrow passage and through a side door, which he opened for us",
		"Within there was a small corridor, which ended in a very massive iron gate.",
		"We were seated at breakfast one morning, my wife and I, when the maid brought in a telegram. It was from Sherlock Holmes and ran in this way",
	}
	for _, s := range sentences {
		if strings.ContainsRune(Censor(s), '*') {
			t.Error("Expected false, got false from sentence", s)
		}
	}
}

func TestSanitize(t *testing.T) {
	expectedString := "whatthefuckisyourproblem"
	sanitizedString, _ := NewProfanityDetector().sanitize("What the fu_ck is y()ur pr0bl3m?", false)
	if sanitizedString != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, sanitizedString)
	}
}
