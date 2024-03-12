package dbfuncs

import (
	"log"
)

func GetPosts(userId string, page int, batchSize int, usersOwnProfile bool) ([]Post, error) {
	offset := (page - 1) * batchSize
	query := `
    SELECT * FROM Posts
    WHERE CreatorId = ?
    ORDER BY CreatedAt DESC
    LIMIT ? OFFSET ?
`
	rows, err := db.Query(query, userId, batchSize, offset)
	if err != nil {
		log.Println(err, "GetPost")
		return nil, err
	}
	defer rows.Close()
	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id,
			&post.Title,
			&post.Body,
			&post.CreatorId,
			&post.CreatedAt,
			&post.Image,
			&post.PrivacyLevel)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// Skip superprivate posts if the viewer does not have access.
		if !usersOwnProfile || post.PrivacyLevel == "superprivate" {
			allowed, err := CheckSuperprivateAccess(post.Id, userId) // Changed from post.Id.String() now that dbfuncs.Post.Id is of type string.
			if err != nil {
				log.Println(err)
				return nil, err
			}
			if !allowed {
				continue
			}
		}

		posts = append(posts, post)
	}

	// If some posts could not be displayed because they were superprivate,
	// recursively call GetPosts till we have 10 posts.
	if !usersOwnProfile {
		if len(posts) < 10 {
			shortfall := 10 - len(posts)
			page++
			morePosts, err := GetPosts(userId, page, shortfall, usersOwnProfile)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			posts = append(posts, morePosts...)
		}
	}

	// Swap order to display latest posts at the bottom.
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}

	return posts, nil
}
