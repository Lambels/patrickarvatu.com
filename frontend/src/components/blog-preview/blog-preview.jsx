function BlogPreview({ blog }) {
  return (
    <a href={`/blog/${blog.title}`}>
      <div className="flex justify-center">
        <div className="flex flex-col lg:min-w-[60rem] lg:min-h-[20rem] md:flex-row md:max-w-xl rounded-lg shadow-lg hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
          <img
            className="w-full lg:min-w-[20rem] md:h-auto object-cover md:w-48 rounded-t-lg md:rounded-none md:rounded-l-lg"
            src={`${process.env.NEXT_PUBLIC_API_URL}/v1/images/blogs/${blog.id}.jpe`}
            alt=""
          />
          <div className="p-6 flex flex-col justify-start">
            <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
              {blog.title}
            </h5>
            <p className="font-normal text-gray-700 dark:text-gray-400">
              {blog.description}
            </p>
            <p className="text-gray-600 text-xs mt-auto">{`Written by Patrick Arvatu At ${blog.updatedAt}`}</p>
          </div>
        </div>
      </div>
    </a>
  );
}

export default BlogPreview;
