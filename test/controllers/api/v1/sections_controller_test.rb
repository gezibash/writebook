require "test_helper"

class Api::V1::SectionsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = users(:david)
    @user.update!(api_token: "test-token-david")
    @headers = { "Authorization" => "Bearer test-token-david" }
    @book = books(:handbook)
  end

  test "index lists sections in a book" do
    post api_v1_book_sections_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "Intro Section" }, section: { body: "overview", theme: "blue" } }

    get api_v1_book_sections_url(@book), headers: @headers, as: :json

    assert_response :success
    json = JSON.parse(response.body)
    assert json.is_a?(Array)
    assert json.any? { |s| s["title"] == "Intro Section" }
  end

  test "create adds a section to a book" do
    assert_difference -> { @book.leaves.count }, +1 do
      post api_v1_book_sections_url(@book), headers: @headers, as: :json,
        params: { leaf: { title: "Part 1" }, section: { body: "First part", theme: "green" } }
    end

    assert_response :created
    json = JSON.parse(response.body)
    assert_equal "Part 1", json["title"]
    assert_equal "Section", json["type"]
    assert_equal "green", json["theme"]
  end

  test "show returns a section" do
    post api_v1_book_sections_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "Visible" }, section: { body: "here", theme: "blue" } }
    leaf_id = JSON.parse(response.body)["id"]

    get api_v1_book_section_url(@book, leaf_id), headers: @headers, as: :json

    assert_response :success
    json = JSON.parse(response.body)
    assert_equal "Visible", json["title"]
    assert_equal "blue", json["theme"]
  end

  test "update modifies a section" do
    post api_v1_book_sections_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "Draft" }, section: { body: "initial" } }
    leaf_id = JSON.parse(response.body)["id"]

    patch api_v1_book_section_url(@book, leaf_id), headers: @headers, as: :json,
      params: { leaf: { title: "Final Section" }, section: { body: "updated", theme: "magenta" } }

    assert_response :success
    json = JSON.parse(response.body)
    assert_equal "Final Section", json["title"]
    assert_equal "magenta", json["theme"]
  end

  test "destroy trashes a section" do
    post api_v1_book_sections_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "To Delete" }, section: { body: "bye" } }
    leaf_id = JSON.parse(response.body)["id"]

    delete api_v1_book_section_url(@book, leaf_id), headers: @headers, as: :json

    assert_response :no_content
    assert_equal "trashed", Leaf.find(leaf_id).status
  end

  test "rejects unauthenticated requests" do
    post api_v1_book_sections_url(@book), as: :json,
      params: { leaf: { title: "Nope" }, section: { body: "nope" } }

    assert_response :unauthorized
  end
end
