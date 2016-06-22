import com.google.gson.Gson
import com.google.gson.GsonBuilder
import opennlp.tools.sentdetect.SentenceDetectorME
import opennlp.tools.sentdetect.SentenceModel
import java.io.FileInputStream
import java.io.FileReader
import java.io.FileWriter
import java.util.*

val sentenceDetector = sentenceDetector()

fun main(args : Array<String>) {
    val gson = GsonBuilder().setPrettyPrinting().disableHtmlEscaping().create();

    val json = fetchRawReviews(gson, "../out/reviews_raw.json")

    val formattedList = ArrayList<String>()

    for (placeDetails in json) {
        if (placeDetails.reviews == null) {
            continue
        }
        val reviews = filterReviews(placeDetails.reviews)
        val formatted = formatReviews(reviews, placeDetails.url)
        formattedList.addAll(formatted)
    }
    Collections.shuffle(formattedList)
    FileWriter("../out/reviews_processed.json").use {
        gson.toJson(formattedList, it)
    }
    print("Done.")
}

private fun fetchRawReviews(gson: Gson, path: String): Array<PlaceDetails> {
    FileReader(path).use {
        return gson.fromJson(it, Array<PlaceDetails>::class.java)
    }
}

fun formatReviews(reviews: List<Review>, placeUrl: String): List<String> {
    val ret = ArrayList<String>();
    for (review in reviews) {
        // TODO: Just truncate at the last possible word boundary instead of continuing.
        val formatted = formatReview(review, placeUrl) ?: continue
        ret.add(formatted)
    }
    return ret;
}

// Tweet length - 23 for URL - 5 for ★★★★★, -2 for separator spaces
const val MAX_REVIEW_TEXT_LENGTH:Int = 140 - 23 - 5 - 2

fun formatReview(review: Review, placeUrl: String): String? {
    var text:String = review.text.trim();
    if (text.length > MAX_REVIEW_TEXT_LENGTH) {
        val shortened = shortenReviewText(text) ?: return null
        text = shortened
    }
    println("""${"★".repeat(review.rating)} $text $placeUrl""")
    return """${"★".repeat(review.rating)} $text $placeUrl"""
}
/*
Logic for review formatting.

1. Detect and shuffle sentences.
2. Add them sequentially to stay under the limit. For each sentence
    - Generate result from adding sentence
    - Check that it's under the limit.
    - Otherwise, revert.
Joining logic:
    - If this is the first sentence, just add the sentence. Done.
    - Sort the indices of all sentences. Add the first sentence. For each extra sentence:
        - Check if the last sentence ends with a full stop. If it does, change that to a "... ".
        - Otherwise, add a "... "
        - Add the sentence.
*/
fun shortenReviewText(review: String): String? {
    val sentences = sentenceDetector.sentDetect(review)
    // shuffle indices
    val indices = listFromTo(0, sentences.size)
    Collections.shuffle(indices)
    var chosen = ArrayList<Int>()
    var currentReviewText:String? = null;

    // Choose these indices until the total is MAX_REVIEW_TEXT_LENGTH chars.
    for (idx in indices) {
        val tryChosen = ArrayList(chosen)
        tryChosen.add(idx)
        val spliced = spliceReview(sentences, tryChosen)
        if (spliced.length == MAX_REVIEW_TEXT_LENGTH) {
            return spliced
        }
        if (spliced.length < MAX_REVIEW_TEXT_LENGTH) {
            chosen = tryChosen
            currentReviewText = spliced
        }
    }
    return currentReviewText
}

fun spliceReview(sentences: Array<String>, chosenSentences: List<Int>): String {
    // Sort the list of chosen indexes, because if there's two in a row, we'll join them using
    // a space instead of an ellipsis and a space.
    val chosen = ArrayList(chosenSentences);
    Collections.sort(chosen)
    var lastIdx = -1;
    var builder = StringBuilder();
    for (idx in chosen) {
        val sentence = sentences[idx]
        when (lastIdx) {
        // Special case. First time, no joining character.
            -1 -> builder.append(sentence)
        // Consecutive sentences
            idx - 1 -> builder.append(" ").append(sentence)
        // Non-consecutive sentences.
            else -> {
                // If the last sentence ended with a period, remove it, because we're about to add
                // an ellipsis. Also trim whitespace. Sometimes people add multiple dots, so sort
                // them out too.
                builder = StringBuilder(builder.trimEnd('.', ' ', '\t', '\n'));
                builder.append("… ").append(sentences[idx])
            }
        }
        lastIdx = idx;
    }
    return builder.toString()
}

fun listFromTo(start: Int, end: Int): ArrayList<Int> {
    val ret = ArrayList<Int>(end-start)
    for (i in start..(end-1)) {
        ret.add(i)
    }
    return ret
}


fun filterReviews(reviews: Array<Review>): List<Review> {
    val ret = ArrayList<Review>()
    for (review in reviews) {
        if (reviewIsUsable(review)) {
            ret.add(review)
        }
    }
    return ret
}

fun reviewIsUsable(review: Review): Boolean {
    return review.language == "en"
        && review.text.trim() != ""
}

fun sentenceDetector(): SentenceDetectorME {
    FileInputStream("../data/en-sent.bin").use {
        val model = SentenceModel(it)
        return SentenceDetectorME(model)
    }
}

data class PlaceDetails(
        @Suppress("ArrayInDataClass") val reviews: Array<Review>?,
        val name: String,
        val url: String)

data class Review(val rating: Int, val text: String, val language: String)
