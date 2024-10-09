package schema

import "github.com/3rd_rec/air_api_tool/consts"

const (
	descriptionCustomField = "CustomField" // Cannot be separated by spaces, otherwise the color display will be abnormal.
)

// The bool type has been deprecated
var supportedFieldType = map[string]bool{
	consts.FieldTypeInt32:      true,
	consts.FieldTypeInt64:      true,
	consts.FieldTypeFloat:      true,
	consts.FieldTypeDouble:     true,
	consts.FieldTypeString:     true,
	consts.FieldTypeJSONString: true,
}

type Field struct {
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	Description  string `json:"description,omitempty"`
	ExampleValue string `json:"example_value,omitempty"`
	Custom       bool   `json:"custom,omitempty"`
}

// industry -> table -> schema
var defaultIndustryTableSchema = map[string]map[string][]*Field{
	consts.IndustrySaasRetail: {
		consts.TableNameUser:      defaultSaasRetailUserSchema,
		consts.TableNameProduct:   defaultSaasRetailProductSchema,
		consts.TableNameUserEvent: defaultSaasRetailUserEventSchema,
	},
	consts.IndustrySaasContent: {
		consts.TableNameUser:      defaultSaasContentUserSchema,
		consts.TableNameContent:   defaultSaasContentContentSchema,
		consts.TableNameUserEvent: defaultSaasContentUserEventSchema,
	},
}

// saas retail default schema configuration
var (
	defaultSaasRetailUserSchema = []*Field{
		{
			Name:         "user_id",
			Type:         "string",
			ExampleValue: "\"1457789\"",
			Description:  "The unique user identifier. Be consistent with the user_id in event table. Will be used as unique identifier in PredictRequest, traffic split of AB experiment etc. Device ID or member ID is often used as user_id here.\nNote: \n1. If you want to encrypt the id and use hashed values here, please ensure hashed id is consistent for the same user. \n2. If your users often switch between login/logout status (In web or mobile application), you might get inconsistent IDs (member v.s. visitor) for the same user. To prevent this, use consistent ID like device ID.",
		},
		{
			Name:         "tags",
			Type:         "json_string",
			ExampleValue: "\"[\\\"new user\\\",\\\"low purchasing power\\\",\\\"bargain seeker\\\"]\"\nor use python code:\njson.dumps([\"new user\", \"low purchasing power\", \"bargain seeker\"])",
			Description:  "The (internal) tags of the given user. Format into JSON serialized string. For example, \"[\\\"new user\\\",\\\"low purchasing power\\\",\\\"bargain seeker\\\"]\".",
		},
		{
			Name:         "language",
			Type:         "string",
			ExampleValue: "\"English\"",
			Description:  "The language(s) set by this user.",
		},
		{
			Name:         "registration_timestamp",
			Type:         "int64",
			ExampleValue: "1623593487",
			Description:  "The timestamp when the given user first activated or registered.",
		},
		{
			Name:         "activation_channel",
			Type:         "string",
			ExampleValue: "\"AppStore\"",
			Description:  "The channel through which user onboarded/signed up. For example, \"AppStore\", \"GoogleAds\", \"FacebookAds\".",
		},
		{
			Name:         "country",
			Type:         "string",
			ExampleValue: "\"USA\"",
			Description:  "The default country set by user, which may differ from their live location.",
		},
		{
			Name:         "province",
			Type:         "string",
			ExampleValue: "\"Texas\"",
			Description:  "The default province set by user, which may differ from their live location.",
		},
		{
			Name:         "city",
			Type:         "string",
			ExampleValue: "\"Kirkland\"",
			Description:  "The default city set by user, which may differ from their live location.",
		},
		{
			Name:         "district",
			Type:         "string",
			ExampleValue: "\"King County\"",
			Description:  "The default district set by user, which may differ from their live location.",
		},
		{
			Name:         "membership_level",
			Type:         "string",
			ExampleValue: "\"silver\"",
			Description:  "The level/tier of the user's membership. For example, \"silver\", \"elite\", \"4\", \"5\".",
		},
		{
			Name:         "gender",
			Type:         "string",
			ExampleValue: "\"male\"",
			Description:  "The gender of the given user. For example, \"male\", \"female\", \"other\".",
		},
		{
			Name:         "age",
			Type:         "string",
			ExampleValue: "\"23\"",
			Description:  "The age of the given user, can be an (estimated) single value, or a range. For example, \"23\", \"18-25\", \"0-15\", \"50-100\".",
		},
	}
	defaultSaasRetailProductSchema = []*Field{
		{
			Name:         "product_id",
			Type:         "string",
			ExampleValue: "\"632461\"",
			Description:  "The unique identifier for the product, can be series_id/entity_id/other unique identifier.",
		},
		{
			Name:         "is_recommendable",
			Type:         "int32",
			ExampleValue: "1",
			Description:  "True (or 1) if the content is recommendable (i.e. to return the content in the recommendation result).\nNote: Even if a content isn't recommendable, include it still as users might have interacted with such content in the past, hence providing insights into behavioural propensities.",
		},
		{
			Name:         "current_price",
			Type:         "float",
			ExampleValue: "49.99",
			Description:  "The current/displayed/discounted price of the product. Round to 2.d.p.",
		},
		{
			Name:         "original_price",
			Type:         "float",
			ExampleValue: "69.98",
			Description:  "The original price of the product. Round to 2.d.p.",
		},
		{
			Name:         "publish_timestamp",
			Type:         "int64",
			ExampleValue: "1623193487",
			Description:  "The timestamp when the product was published.",
		},
		{
			Name:         "categories",
			Type:         "json_string",
			ExampleValue: "\"[{\\\"category_depth\\\":1,\\\"category_nodes\\\":[{\\\"id_or_name\\\":\\\"Shoes\\\"}]},{\\\"category_depth\\\":2,\\\"category_nodes\\\":[{\\\"id_or_name\\\":\\\"Men's Shoes\\\"}]}]\"\nor use python code:\njson.dumps([{\"category_depth\": 1,\"category_nodes\": [{\"id_or_name\": \"Shoes\"}]},{\"category_depth\": 2,\"category_nodes\": [{\"id_or_name\": \"Men Shoes\"}]}])",
			Description:  "The (sub)categories the content fall under. Format requirements: \n1. JSON serialised string\n2. Depth starts from 1, in consecutive postive integers\n3. \"Category_nodes\" should not contain empty value. If empty value exists, replace with \"null\" . \n4. Only one \"id_or_name\" key-value pair is allowed under each \"category_nodes\"",
		},
		{
			Name:         "tags",
			Type:         "json_string",
			ExampleValue: "\"[\\\"New Product\\\",\\\"Summer Product \\\"]\"\nor use python code:\njson.dumps([\"New Product\",\"Summer Product\"])",
			Description:  "The (internal) label of the product. Format into JSON serialized string. ",
		},
		{
			Name:         "title",
			Type:         "string",
			ExampleValue: "adidas Men's Yeezy Boost 350 V2 Grey/Borang/Dgsogr",
			Description:  "The title/name of the product.",
		},
		{
			Name:         "brands",
			Type:         "string",
			ExampleValue: "\"Adidas\"",
			Description:  "The brand of the product.",
		},
		{
			Name:         "user_rating",
			Type:         "float",
			ExampleValue: "0.25",
			Description:  "The rating given by the users. Non-negative value.",
		},
		{
			Name:         "sold_count",
			Type:         "int32",
			ExampleValue: "60",
			Description:  "The number of sales of the content if it is paid content.",
		},
		{
			Name:         "product_group_id",
			Type:         "string",
			ExampleValue: "\"1356\"",
			Description:  "The ID of the group/unit for products with common characteristics.",
		},
		{
			Name:         "display_cover_multimedia_url",
			Type:         "json_string",
			ExampleValue: "\"[\\\"https://images-na.ssl-images-amazon.com/images/I/81WmojBxvbL._AC_UL1500.jpg\\\"]\"\nor use python code:\njson.dumps([\"https://images-na.ssl-images-amazon.com/images/I/81WmojBxvbL._AC_UL1500.jpg\"])",
			Description:  "The URL of the cover multimedia for the product. Format into JSON serialized string.",
		},
		{
			Name:         "comment_count",
			Type:         "int32",
			ExampleValue: "100",
			Description:  "The number of comments of the content. Non-negative value.",
		},
		{
			Name:         "source",
			Type:         "string",
			ExampleValue: "\"self\"",
			Description:  "The source of the product.",
		},
		{
			Name:         "seller_id",
			Type:         "string",
			ExampleValue: "\"43485\"",
			Description:  "The ID of the seller.",
		},
		{
			Name:         "seller_level",
			Type:         "string",
			ExampleValue: "\"1\"",
			Description:  "The tier/level of the seller.",
		},
		{
			Name:         "seller_rating",
			Type:         "float",
			ExampleValue: "3.5",
			Description:  "The seller's rating given by the customers. Non-negative value.",
		},
	}
	defaultSaasRetailUserEventSchema = []*Field{
		{
			Name:         "product_id",
			Type:         "string",
			ExampleValue: "\"1457789\"",
			Description:  "The unique user identifier. Be consistent with the user_id and content_owner_id in user and content table.",
		},
		{
			Name:         "event_type",
			Type:         "string",
			ExampleValue: "\"purchase\"",
			Description:  "The user event type. Predefined values are: \"impression\", \"click\", \"add-to-cart\", \"remove-from-cart\", \"add-to-favorites\", \"remove-from-favorites\", \"purchase\", \"search\", \"checkout\".",
		},
		{
			Name:         "event_timestamp",
			Type:         "int64",
			ExampleValue: "1640657087",
			Description:  "The Unix timestamp when the event took place.",
		},
		{
			Name:         "scene_name",
			Type:         "string",
			ExampleValue: "\"product detail page\"",
			Description:  "The unique identifier of the (sub)scene where the event took place. Required for \"impression\" and \"click\" events. Be as specific as possible.",
		},
		{
			Name:         "product_id",
			Type:         "string",
			ExampleValue: "\"632461\"",
			Description:  "The ID of the content in respect of the event. Not required when event_type is \"search\". Otherwise it is required. (refer to the Accepted values for event_type above)",
		},
		{
			Name:         "purchase_count",
			Type:         "int32",
			ExampleValue: "20",
			Description:  "The final amount in the purchase record. Required for \"purchase\" event type.",
		},
		{
			Name:         "paid_price",
			Type:         "float",
			ExampleValue: "12.23",
			Description:  "The final amount paid by the user. Required for \"purchase\" event type. Match the value with the \"currency\" used.",
		},
		{
			Name:         "parent_product_id",
			Type:         "string",
			ExampleValue: "\"441356\"",
			Description:  "The ID of the root product in respect of the event. Recommend to provide if user viewed/clicked on the product from a product detail page.",
		},
		{
			Name:         "query",
			Type:         "string",
			ExampleValue: "\"iPad\"",
			Description:  "The search query. Recommend to provide for \"search\" , and \"impression\", \"click\" events ensued from search.",
		},
		{
			Name:         "page_number",
			Type:         "int32",
			ExampleValue: "2",
			Description:  "The page number (of the product) where the event took place. For example: (X = page_number)\n1. Users swipe X times to like the product.\n2. Users go to page X to view the product.\nwhen X=0 (default value), it means all content fits under one page.",
		},
		{
			Name:         "offset",
			Type:         "int32",
			ExampleValue: "10",
			Description:  "The position (of the product) in the specified page where the event took place. Start from 1.",
		},
		{
			Name:         "platform",
			Type:         "string",
			ExampleValue: "\"app\"",
			Description:  "The user's platform. For example, \"app\", \"desktop_web\", \"mobile_web\", \"other\".",
		},
		{
			Name:         "os_type",
			Type:         "string",
			ExampleValue: "\"android\"",
			Description:  "The user's operating system.",
		},
		{
			Name:         "os_version",
			Type:         "string",
			ExampleValue: "\"10\"",
			Description:  "The version of user's operating system.",
		},
		{
			Name:         "app_version",
			Type:         "string",
			ExampleValue: "\"1.0.1\"",
			Description:  "The version of user's application.",
		},
		{
			Name:         "network",
			Type:         "string",
			ExampleValue: "\"3g\"",
			Description:  "The network user used.",
		},
		{
			Name:         "device_model",
			Type:         "string",
			ExampleValue: "\"huawei-mate30\"",
			Description:  "The user's device model.",
		},
		{
			Name:         "attribution_token",
			Type:         "string",
			ExampleValue: "\"eyJpc3MiOiJuaW5naGFvLm5ldCIsImV4cCI6IjE0Mzg5NTU0NDUiLCJuYW1lIjoid2FuZ2hhbyIsImFkbWluIjp0cnVlfQ\"",
			Description:  "Provided by BytePlus. The identifier assigned to all attributed events in the same user session.",
		},
		{
			Name:         "traffic_source",
			Type:         "string",
			ExampleValue: "\"self\"",
			Description:  "The traffic source of the event. To provide upon A/B experiment go-live. Acceptable values are \"self\", \"byteplus\", \"other\"\n\"self\" : from the user's own server\n\"byteplus\": from byteplus's server\n\"other\" : from a third-party server",
		},
		{
			Name:         "currency",
			Type:         "string",
			ExampleValue: "\"USD\"",
			Description:  "The currency used in transaction record. Recommend to provide for \"purchase\" event type. Standardize to USD across multiple countries/regions.",
		},
		{
			Name:         "country",
			Type:         "string",
			ExampleValue: "\"USA\"",
			Description:  "The country where the event took place.",
		},
		{
			Name:         "province",
			Type:         "string",
			ExampleValue: "\"Texas\"",
			Description:  "The province where the event took place.",
		},
		{
			Name:         "city",
			Type:         "string",
			ExampleValue: "\"Kirkland\"",
			Description:  "The city where the event took place.",
		},
		{
			Name:         "district",
			Type:         "string",
			ExampleValue: "\"King County\"",
			Description:  "The district where the event took place.",
		},
	}
)

// saas content default schema configuration
var (
	defaultSaasContentUserSchema = []*Field{
		{
			Name:         "user_id",
			Type:         "string",
			ExampleValue: "\"1457789\"",
			Description:  "The unique user identifier. Be consistent with the user_id and content_owner_id in event and content table. Will be used as unique identifier in PredictRequest, traffic split of AB experiment etc. Device ID or member ID is often used as user_id here.\n\nNote: \n1. If you want to encrypt the id and use hashed values here, please ensure hashed id is consistent for the same user. \n2. If your users often switch between login/logout status (In web or mobile application), you might get inconsistent IDs (member v.s. visitor) for the same user. To prevent this, use consistent ID like device ID.",
		},
		{
			Name:         "tags",
			Type:         "json_string",
			ExampleValue: "\"[\\\"new user\\\",\\\"low purchasing power\\\",\\\"bargain seeker\\\"]\"\nor use python code:\njson.dumps([\"new user\", \"low purchasing power\", \"bargain seeker\"])",
			Description:  "The (internal) tags of the given user. Format into JSON serialized string. For example, \"[\\\"new user\\\",\\\"low purchasing power\\\",\\\"bargain seeker\\\"]\".",
		},
		{
			Name:         "language",
			Type:         "string",
			ExampleValue: "\"English\"",
			Description:  "The language(s) set by this user.",
		},
		{
			Name:         "registration_timestamp",
			Type:         "int64",
			ExampleValue: "1623593487",
			Description:  "The timestamp when the given user first activated or registered.",
		},
		{
			Name:         "subscriber_type",
			Type:         "string",
			ExampleValue: "\"free\"",
			Description:  "The user's subscription tier, whether the user is free/a subscriber.",
		},
		{
			Name:         "membership_level",
			Type:         "string",
			ExampleValue: "\"silver\"",
			Description:  "The level/tier of the user's membership. For example, \"silver\", \"elite\", \"4\", \"5\".",
		},
		{
			Name:         "country",
			Type:         "string",
			ExampleValue: "\"USA\"",
			Description:  "The default country set by user, which may differ from their live location.",
		},
		{
			Name:         "province",
			Type:         "string",
			ExampleValue: "\"Texas\"",
			Description:  "The default province set by user, which may differ from their live location.",
		},
		{
			Name:         "city",
			Type:         "string",
			ExampleValue: "\"Kirkland\"",
			Description:  "The default city set by user, which may differ from their live location.",
		},
		{
			Name:         "district",
			Type:         "string",
			ExampleValue: "\"King County\"",
			Description:  "The default district set by user, which may differ from their live location.",
		},
		{
			Name:         "gender",
			Type:         "string",
			ExampleValue: "\"male\"",
			Description:  "The gender of the given user. For example, \"male\", \"female\", \"other\".",
		},
		{
			Name:         "age",
			Type:         "string",
			ExampleValue: "\"23\"",
			Description:  "The age of the given user, can be an (estimated) single value, or a range. For example, \"23\", \"18-25\", \"0-15\", \"50-100\".",
		},
	}
	defaultSaasContentContentSchema = []*Field{
		{
			Name:         "content_id",
			Type:         "string",
			ExampleValue: "\"632461\"",
			Description:  "The unique identifier for the content, can be series_id/entity_id/video_id/other unique identifier.",
		},
		{
			Name:         "content_type",
			Type:         "string",
			ExampleValue: "\"video\"",
			Description:  "The type of the content.",
		},
		{
			Name:         "is_recommendable",
			Type:         "int32",
			ExampleValue: "1",
			Description:  "True (or 1) if the content is recommendable (i.e. to return the content in the recommendation result).\nNote: Even if a content isn't recommendable, include it still as users might have interacted with such content in the past, hence providing insights into behavioural propensities.",
		},
		{
			Name:         "language",
			Type:         "string",
			ExampleValue: "\"English\"",
			Description:  "The languages used in the content.",
		},
		{
			Name:         "content_title",
			Type:         "string",
			ExampleValue: "\"Green Book Movie Explanation\"",
			Description:  "The title of the content.",
		},
		{
			Name:         "categories",
			Type:         "json_string",
			ExampleValue: "\"[{\\\"category_depth\\\":1,\\\"category_nodes\\\":[{\\\"id_or_name\\\":\\\"Movie\\\"}]},{\\\"category_depth\\\":2,\\\"category_nodes\\\":[{\\\"id_or_name\\\":\\\"Comedy\\\"}]}]\"\nor use python code:\njson.dumps([{\"category_depth\": 1,\"category_nodes\": [{\"id_or_name\": \"Movie\"}]},{\"category_depth\": 2,\"category_nodes\": [{\"id_or_name\": \"Comedy\"}]}])",
			Description:  "The (sub)categories the content fall under. Format requirements: \n1. JSON serialised string\n2. Depth starts from 1, in consecutive postive integers\n3. \"Category_nodes\" should not contain empty value. If empty value exists, replace with \"null\" . \n4. Only one \"id_or_name\" key-value pair is allowed under each \"category_nodes\"",
		},
		{
			Name:         "tags",
			Type:         "json_string",
			ExampleValue: "\"[\\\"New\\\",\\\"Trending\\\"]\"\nor use python code:\njson.dumps([\"New\",\"Trending\"])",
			Description:  "The (internal) label of the content. Format into JSON serialized string.",
		},
		{
			Name:         "collection_id",
			Type:         "string",
			ExampleValue: "\"1342\"",
			Description:  "The ID of the colllection, if any, to which the content belongs.",
		},
		{
			Name:         "publish_timestamp",
			Type:         "int64",
			ExampleValue: "1660035734",
			Description:  "The timestamp when the content was published.",
		},
		{
			Name:         "content_owner_id",
			Type:         "string",
			ExampleValue: "\"1457789\"",
			Description:  "The user_id of the content creator.",
		},
		{
			Name:         "video_duration",
			Type:         "int32",
			ExampleValue: "1200000",
			Description:  "The length of the video in milliseconds.",
		},
		{
			Name:         "description",
			Type:         "string",
			ExampleValue: "\"A brief summary of the main content of the Green Book movie\"",
			Description:  "The detailed description of the content.",
		},
		{
			Name:         "linked_product_id",
			Type:         "json_string",
			ExampleValue: "\"[\\\"632462\\\",\\\"632463\\\"]\"\nor use python code:\njson.dumps([\"632462\", \"632463\"])",
			Description:  "The product_id of the goods sold via the content.",
		},
		{
			Name:         "video_urls",
			Type:         "json_string",
			ExampleValue: "\"[\\\"https://test_video.mov\\\"]\"\nor use python code:\njson.dumps([\"https://test_video.mov\"])",
			Description:  "The URL of the video. Format into JSON serialized string.",
		},
		{
			Name:         "image_urls",
			Type:         "json_string",
			ExampleValue: "\"[\\\"https://images-na.ssl-images-amazon.com/images/I/81WmojBxvbL._AC_UL1500.jpg\\\"]\"\nor use python code:\njson.dumps([\"https://images-na.ssl-images-amazon.com/images/I/81WmojBxvbL._AC_UL1500.jpg\"])",
			Description:  "The URL of the image (e.g. thumbnail). Format into JSON serialized string.",
		},
		{
			Name:         "user_rating",
			Type:         "float",
			ExampleValue: "4.9",
			Description:  "The rating of the content. Non-negative value.",
		},
		{
			Name:         "is_paid_content",
			Type:         "bool",
			ExampleValue: "true",
			Description:  "True (or 1) if the content requires payment/subscription to view.",
		},
		{
			Name:         "source",
			Type:         "string",
			ExampleValue: "\"self\"",
			Description:  "The source of the content.",
		},
		{
			Name:         "current_price",
			Type:         "float",
			ExampleValue: "1300.12",
			Description:  "The price (after discount) of the content in cents if it is paid content.",
		},
		{
			Name:         "original_price",
			Type:         "float",
			ExampleValue: "1600.12",
			Description:  "The price (before discount) of the content in cents if it is paid content.",
		},
	}
	defaultSaasContentUserEventSchema = []*Field{
		{
			Name:         "user_id",
			Type:         "string",
			ExampleValue: "\"1457787\"",
			Description:  "The unique user identifier. Be consistent with the user_id and content_owner_id in user and content table.",
		},
		{
			Name:         "content_id",
			Type:         "string",
			ExampleValue: "\"632461\"",
			Description:  "The ID of the content w.r.t the event. Not required when event_type is \"search\" or \"follow\". Otherwise it is required. (refer to the Accepted values for event_type above)",
		},
		{
			Name:         "event_type",
			Type:         "string",
			ExampleValue: "\"impression\"",
			Description:  "The user event type. Predefined values are: \"impression\", \"click\", \"stay\", \"like\", \"share\", \"comment\", \"follow\", \"favorite\", \"search\", \"cart\", \"checkout\", \"purchase\"",
		},
		{
			Name:         "event_timestamp",
			Type:         "int64",
			ExampleValue: "1640657087",
			Description:  "The Unix timestamp when the event took place.",
		},
		{
			Name:         "scene_name",
			Type:         "string",
			ExampleValue: "\"Home Page\"",
			Description:  "The unique identifier of the (sub)scene where the event took place. Required for \"impression\", \"click\" and \"stay\" events. Be as specific as possible.",
		},
		{
			Name:         "stay_duration",
			Type:         "int32",
			ExampleValue: "150000",
			Description:  "The times user spent viewing the content in millisecond. Required for \"stay\" event.",
		},
		{
			Name:         "parent_content_id",
			Type:         "string",
			ExampleValue: "\"632431\"",
			Description:  "The ID of the root/main content in respect of the event. Recommend to provide if user viewed/clicked on the content from a content detail page.",
		},
		{
			Name:         "query",
			Type:         "string",
			ExampleValue: "\"comedy\"",
			Description:  "The search query. Recommend to provide for \"search\" , and \"impression\", \"click\" events ensued from search.",
		},
		{
			Name:         "content_owner_id",
			Type:         "string",
			ExampleValue: "\"1457789\"",
			Description:  "The user_id of the content owner in respect of the event. Recommend to provide for \"follow\" event.",
		},
		{
			Name:         "purchase_count",
			Type:         "int32",
			ExampleValue: "20",
			Description:  "The final amount in the purchase record. Required for \"purchase\" event.",
		},
		{
			Name:         "paid_price",
			Type:         "float",
			ExampleValue: "12.23",
			Description:  "The final amount paid by the user. Required for \"purchase\" event. Match the value with the \"currency\" used.",
		},
		{
			Name:         "page_number",
			Type:         "int32",
			ExampleValue: "2",
			Description:  "The page number (of the content) where the event took place. For example: (X = page_number)\n1. Users swipe X times to like the content.\n2. Users go to page X to view the content.\nwhen X=0 (default value), it means all content fits under one page.",
		},
		{
			Name:         "offset",
			Type:         "int32",
			ExampleValue: "10",
			Description:  "The position (of the content) in the specified page where the event took place. Start from 1.",
		},
		{
			Name:         "platform",
			Type:         "string",
			ExampleValue: "\"app\"",
			Description:  "The user's platform. For example, \"app\", \"desktop_web\", \"mobile_web\", \"other\".",
		},
		{
			Name:         "os_type",
			Type:         "string",
			ExampleValue: "\"android\"",
			Description:  "The user's operating system.",
		},
		{
			Name:         "os_version",
			Type:         "string",
			ExampleValue: "\"10\"",
			Description:  "The version of user's operating system.",
		},
		{
			Name:         "app_version",
			Type:         "string",
			ExampleValue: "\"1.0.1\"",
			Description:  "The version of user's application.",
		},
		{
			Name:         "network",
			Type:         "string",
			ExampleValue: "\"3g\"",
			Description:  "The network user used.",
		},
		{
			Name:         "device_model",
			Type:         "string",
			ExampleValue: "\"huawei-mate30\"",
			Description:  "The user's device model.",
		},
		{
			Name:         "attribution_token",
			Type:         "string",
			ExampleValue: "\"eyJpc3MiOiJuaW5naGFvLm5ldCIsImV4cCI6IjE0Mzg5NTU0NDUiLCJuYW1lIjoid2FuZ2hhbyIsImFkbWluIjp0cnVlfQ\"",
			Description:  "Provided by BytePlus. The identifier assigned to all attributed events in the same user session.",
		},
		{
			Name:         "traffic_source",
			Type:         "string",
			ExampleValue: "\"self\"",
			Description:  "The traffic source of the event. To provide upon A/B experiment go-live. Acceptable values are \"self\", \"byteplus\", \"other\"\n\"self\" : from the user's own server\n\"byteplus\": from byteplus's server\n\"other\" : from a third-party server",
		},
		{
			Name:         "currency",
			Type:         "string",
			ExampleValue: "\"USD\"",
			Description:  "The currency used in transaction record. Recommend to provide for \"purchase\" event type. Standardize to USD across multiple countries/regions.",
		},
		{
			Name:         "country",
			Type:         "string",
			ExampleValue: "\"USA\"",
			Description:  "The country where the event took place.",
		},
		{
			Name:         "province",
			Type:         "string",
			ExampleValue: "\"Texas\"",
			Description:  "The province where the event took place.",
		},
		{
			Name:         "city",
			Type:         "string",
			ExampleValue: "\"Kirkland\"",
			Description:  "The city where the event took place.",
		},
		{
			Name:         "district",
			Type:         "string",
			ExampleValue: "\"King County\"",
			Description:  "The district where the event took place.",
		},
	}
)
