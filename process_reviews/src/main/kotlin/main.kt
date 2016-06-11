import com.google.gson.Gson;
import java.io.FileReader
import java.util.*

fun main(args : Array<String>) {
    val reader = FileReader("../out/reviews_raw.json")
    val json = Gson().fromJson(reader, Array<PlaceDetails>::class.java);

    val formattedList = ArrayList<String>()

    for (placeDetails in json) {
        if (placeDetails.reviews == null) {
            continue
        }
        val reviews = filterReviews(placeDetails.reviews)
        val formatted = formatReviews(reviews, placeDetails.url)
        formattedList.addAll(formatted)
    }

}

fun formatReviews(reviews: List<Review>, placeUrl: String): List<String> {
    val ret = ArrayList<String>();
    for (review in reviews) {
        ret.add(formatReview(review, placeUrl))
    }
    println(ret)
    return ret;
}

fun formatReview(review: Review, placeUrl: String): String {
    val text = review.text;
    return """${"★".repeat(review.rating)} $text $placeUrl"""
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

data class PlaceDetails(val reviews: Array<Review>?, val name: String, val url: String)

data class Review(val rating: Int, val text: String, val language: String)
