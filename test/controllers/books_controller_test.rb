require "test_helper"

class BooksControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in :kevin
  end

  test "index lists the current user's books" do
    get root_url

    assert_response :success
    assert_select "h2", text: "Handbook"
    assert_select "h2", text: "Manual", count: 0
  end

  test "index includes published books, even when the user does not have access" do
    books(:manual).update!(published: true)

    get root_url

    assert_response :success
    assert_select "h2", text: "Handbook"
    assert_select "h2", text: "Manual"
  end

  test "index shows published books when not logged in" do
    books(:manual).update!(published: true)

    sign_out
    get root_url

    assert_response :success
    assert_select "h2", text: "Handbook", count: 0
    assert_select "h2", text: "Manual"
  end

  test "index redirects to login if not signed in and no published books exist" do
    sign_out
    get root_url

    assert_redirected_to new_session_url
  end

  test "new renders cover style picker" do
    get new_book_url

    assert_response :success
    assert_select "[data-controller='book-cover-preview']"
    assert_select "input[name='book[cover_seed]'][type='hidden'][value]", count: 1
    assert_select "[data-action='book-cover-preview#refreshSeed']"
    assert_select "input[name='book[cover_style]'][value='glass'][checked='checked']"
    assert_select "input[name='book[cover_style]'][value='rings']"
  end

  test "create makes the current user an editor" do
    cover_seed = "feedfacecafebeef"

    assert_difference -> { Book.count }, +1 do
      post books_url, params: { book: { title: "New Book", everyone_access: false, cover_style: "glass", cover_seed: cover_seed } }
    end

    assert_redirected_to book_slug_url(Book.last)

    book = Book.last
    assert_equal "New Book", book.title
    assert_equal "glass", book.cover_style
    assert_equal cover_seed, book.cover_seed
    assert_equal 1, Book.last.accesses.count

    assert book.editable?(user: users(:kevin))
  end

  test "create sets additional accesses" do
    sign_in :jason
    assert_difference -> { Book.count }, +1 do
      post books_url, params: { book: { title: "New Book", everyone_access: false }, "editor_ids[]": users(:jz).id, "reader_ids[]": users(:kevin).id }
    end

    book = Book.last
    assert_equal "New Book", book.title
    assert_equal 3, Book.last.accesses.count

    assert book.editable?(user: users(:jz))

    assert book.accessable?(user: users(:kevin))
    assert_not book.editable?(user: users(:kevin))
  end

  test "show only shows books the current user can access" do
    get book_slug_url(books(:manual))
    assert_response :not_found

    get book_slug_url(books(:handbook))
    assert_response :success
  end

  test "show includes OG metadata for public access" do
    get book_slug_url(books(:handbook))
    assert_response :success

    assert_select "meta[property='og:title'][content='Handbook']"
    assert_select "meta[property='og:url'][content='#{book_slug_url(books(:handbook))}']"
  end
end
