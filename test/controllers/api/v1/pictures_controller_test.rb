require "test_helper"

class Api::V1::PicturesControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = users(:david)
    @user.update!(api_token: "test-token-david")
    @headers = { "Authorization" => "Bearer test-token-david" }
    @book = books(:handbook)
  end

  test "index lists pictures in a book" do
    get api_v1_book_pictures_url(@book), headers: @headers, as: :json

    assert_response :success
    json = JSON.parse(response.body)
    assert json.is_a?(Array)
  end

  test "create adds a picture with image to a book" do
    image = fixture_file_upload("white-rabbit.webp", "image/webp")

    assert_difference -> { @book.leaves.count }, +1 do
      post api_v1_book_pictures_url(@book), headers: @headers,
        params: { leaf: { title: "Rabbit Photo" }, picture: { image: image, caption: "A white rabbit" } }
    end

    assert_response :created
    json = JSON.parse(response.body)
    assert_equal "Rabbit Photo", json["title"]
    assert_equal "Picture", json["type"]
    assert_equal "A white rabbit", json["caption"]
    assert json["has_image"]
    assert json["image_url"].present?
  end

  test "show returns a picture" do
    image = fixture_file_upload("reading.webp", "image/webp")
    post api_v1_book_pictures_url(@book), headers: @headers,
      params: { leaf: { title: "Reading" }, picture: { image: image, caption: "Reading time" } }
    leaf_id = JSON.parse(response.body)["id"]

    get api_v1_book_picture_url(@book, leaf_id), headers: @headers, as: :json

    assert_response :success
    json = JSON.parse(response.body)
    assert_equal "Reading", json["title"]
    assert_equal "Reading time", json["caption"]
    assert json["has_image"]
  end

  test "update modifies a picture" do
    image = fixture_file_upload("white-rabbit.webp", "image/webp")
    post api_v1_book_pictures_url(@book), headers: @headers,
      params: { leaf: { title: "Draft" }, picture: { image: image } }
    leaf_id = JSON.parse(response.body)["id"]

    patch api_v1_book_picture_url(@book, leaf_id), headers: @headers,
      params: { leaf: { title: "Final Photo" }, picture: { caption: "Updated caption" } }

    assert_response :success
    json = JSON.parse(response.body)
    assert_equal "Final Photo", json["title"]
    assert_equal "Updated caption", json["caption"]
  end

  test "destroy trashes a picture" do
    image = fixture_file_upload("white-rabbit.webp", "image/webp")
    post api_v1_book_pictures_url(@book), headers: @headers,
      params: { leaf: { title: "To Delete" }, picture: { image: image } }
    leaf_id = JSON.parse(response.body)["id"]

    delete api_v1_book_picture_url(@book, leaf_id), headers: @headers, as: :json

    assert_response :no_content
    assert_equal "trashed", Leaf.find(leaf_id).status
  end

  test "rejects unauthenticated requests" do
    get api_v1_book_pictures_url(@book), as: :json

    assert_response :unauthorized
  end
end
