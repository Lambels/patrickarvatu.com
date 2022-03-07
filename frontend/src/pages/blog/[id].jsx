import { NextSeo } from "next-seo";
import Blog from "../../components/blog";

function BlogPage({ data }) {
  return (
    <>
      <NextSeo title={data.title} />
      <Blog blog={data}/>
    </>
  );
}

export async function getServerSideProps({ query }) {
  const response = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/v1/blogs/${query.id}`
  );

  const data = await response.json();
  return { props: { data } };
}

export default BlogPage;
