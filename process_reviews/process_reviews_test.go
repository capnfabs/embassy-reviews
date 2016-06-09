package main

import "testing"

const testData = `I was at the Nigeria Embassy Berlin on the 21th of Jan. 2014. I got there at about 10:30am and I got my passport about 11:30 am.
I got the best service ever. The staffs there are very nice and friendly.
I thank you all for a Job well done!!!!!!.
Keep it up!!!!!.`

func TestLimitChooseSentence(t *testing.T) {
	resp := limitChooseSentence(testData, 110)
	if len(resp) > 110 {
		t.Errorf("Got a string longer than expected:\n%s", resp)
	}
}
