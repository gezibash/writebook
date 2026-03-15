module Api
  module V1
    class BooksController < BaseController
      before_action :set_book, only: %i[ show update destroy ]
      before_action -> { ensure_editable(@book) }, only: %i[ update destroy ]

      def index
        books = Book.accessable_or_published.ordered
        render json: books.map { |b| book_json(b) }
      end

      def show
        leaves = @book.leaves.active.with_leafables.positioned
        render json: book_json(@book).merge(
          leaves: leaves.map { |l| leaf_json(l) }
        )
      end

      def create
        book = Book.create!(book_params)

        editors = [ Current.user.id ]
        readers = [ Current.user.id ]
        book.update_access(editors: editors, readers: readers)

        render json: book_json(book), status: :created
      end

      def update
        @book.update!(book_params)
        render json: book_json(@book)
      end

      def destroy
        @book.destroy
        head :no_content
      end

      private
        def set_book
          @book = Book.accessable_or_published.find(params[:id])
        end

        def book_params
          params.require(:book).permit(:title, :subtitle, :author, :everyone_access, :theme)
        end

        def book_json(book)
          {
            id: book.id,
            title: book.title,
            subtitle: book.subtitle,
            author: book.author,
            slug: book.slug,
            published: book.published,
            theme: book.theme,
            everyone_access: book.everyone_access,
            created_at: book.created_at,
            updated_at: book.updated_at
          }
        end

        def leaf_json(leaf)
          json = {
            id: leaf.id,
            title: leaf.title,
            type: leaf.leafable_type,
            status: leaf.status,
            position: leaf.position_score,
            slug: leaf.slug,
            created_at: leaf.created_at,
            updated_at: leaf.updated_at
          }

          case leaf.leafable
          when Page
            json[:body] = leaf.leafable.body&.content.to_s
          when Section
            json[:body] = leaf.leafable.body
            json[:theme] = leaf.leafable.theme
          when Picture
            json[:caption] = leaf.leafable.caption
            json[:has_image] = leaf.leafable.image.attached?
          end

          json
        end
    end
  end
end
