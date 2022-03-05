import { NextSeo } from "next-seo";
import AdminBlogPreview from "../../components/admin-blog-preview/admin-blog-preview";
import BlogPreview from "../../components/blog-preview/blog-preview";
import BlogsPreviewLayout from "../../components/blogs-preview-layout";
import { useAuth } from "../../store/auth-context";

function Index({ blogData }) {
  const { isAdmin } = useAuth();

  if (isAdmin) {
    return (
      <>
        <NextSeo title="Blogs" />
        <BlogsPreviewLayout>
          {blogData.blogs?.map((obj, i) => {
            return <AdminBlogPreview blog={obj} key={i} />;
          })}
        </BlogsPreviewLayout>
      </>
    );
  } else {
    return (
      <>
        <NextSeo title="Blogs" />
        <BlogsPreviewLayout>
          {blogData.blogs?.map((obj, i) => {
              return <BlogPreview blog={obj} key={i} />;
          })}
        </BlogsPreviewLayout>
      </>
    );
  }
}

export async function getServerSideProps() {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/v1/blogs`);

  const blogData = await response.json();
  return { props: { blogData } };
}

export default Index;
