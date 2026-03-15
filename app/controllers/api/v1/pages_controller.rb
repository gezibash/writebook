module Api
  module V1
    class PagesController < BaseController
      before_action :set_book
      before_action -> { ensure_editable(@book) }, except: %i[ index ]
      before_action :set_leaf, only: %i[ show update destroy ]

      def index
        pages = @book.leaves.active.with_leafables.positioned.where(leafable_type: "Page")
        render json: pages.map { |l| page_json(l) }
      end

      def show
        render json: page_json(@leaf)
      end

      def create
        page = Page.new(page_params)
        leaf = @book.press(page, leaf_params)

        if position = params[:position]&.to_i
          leaf.move_to_position(position)
        end

        render json: page_json(leaf), status: :created
      end

      def update
        @leaf.edit(leafable_params: page_params, leaf_params: leaf_params)
        render json: page_json(@leaf.reload)
      end

      def destroy
        @leaf.trashed!
        head :no_content
      end

      private
        def set_book
          @book = Book.accessable_or_published.find(params[:book_id])
        end

        def set_leaf
          @leaf = @book.leaves.active.find(params[:id])
        end

        def leaf_params
          params.fetch(:leaf, {}).permit(:title).reverse_merge(title: "Untitled")
        end

        def page_params
          params.fetch(:page, {}).permit(:body)
        end

        def page_json(leaf)
          {
            id: leaf.id,
            title: leaf.title,
            type: "Page",
            body: leaf.leafable.body&.content.to_s,
            status: leaf.status,
            position: leaf.position_score,
            slug: leaf.slug,
            created_at: leaf.created_at,
            updated_at: leaf.updated_at
          }
        end
    end
  end
end
