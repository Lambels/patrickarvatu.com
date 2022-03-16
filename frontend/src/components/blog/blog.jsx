import ReactMarkdown from "react-markdown";
import { dracula } from "react-syntax-highlighter/dist/cjs/styles/prism";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import React from "react";
import remarkGfm from "remark-gfm";

function Blog({ blog }) {
  return (
    <div className="justify-center">
      <article title={blog.title} className="text-white pt-3">
        {blog.subBlogs?.map((obj, i) => {
          return (
            <section title={obj.title} key={i}>
              <a
                href={`#${obj.title}`}
                className="no-underline decoration-2 hover:underline hover:decoration-indigo-500"
              >
                <h2 className="text-6xl font-medium px-10 py-10" id={obj.title}>
                  {obj.title}
                </h2>
              </a>
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                components={{
                  h1({ node, ...props }) {
                    return (
                      <h1
                        className="text-3xl font-sans text-gray-100 underline decoration-slate-600 px-10 py-5"
                        {...props}
                      />
                    );
                  },
                  h2({ node, ...props }) {
                    return (
                      <h2
                        className="text-2xl font-sans px-10 py-5"
                        {...props}
                      />
                    );
                  },
                  h3({ node, ...props }) {
                    return (
                      <h3 className="text-xl font-sans px-10 py-5" {...props} />
                    );
                  },
                  p({ node, ...props }) {
                    return (
                      <p
                        className="text-lg font-light text-gray-300 px-10 pb-2"
                        {...props}
                      />
                    );
                  },
                  a({ node, ...props }) {
                    return (
                      <a
                        className="text-lg font-bold text-white underline decoration-indigo-500 decoration-2 hover:decoration-[3px]"
                        {...props}
                      />
                    );
                  },
                  img({ node, ...props }) {
                    return (
                      <div className="grid place-items-center p-5">
                        <img {...props} />
                      </div>
                    )
                  },
                  blockquote({ node, ...props }) {
                    return (
                      <div className="px-10 py-3">
                      <blockquote className="italic border-l-4 border-gray-700 py-3">
                        <p className="text-lg" {...props} />
                      </blockquote>
                      </div>
                    );
                  },
                  code({ node, inline, className, children, ...props }) {
                    const match = /language-(\w+)/.exec(className || "");

                    return !inline && match ? (
                      <div className="px-10 py-5">
                        <SyntaxHighlighter
                          language={match[1]}
                          PreTag="div"
                          style={dracula}
                          showLineNumbers
                          {...props}
                        >
                          {String(children).replace(/\n$/, "")}
                        </SyntaxHighlighter>
                      </div>
                    ) : (
                      <code className={className} {...props}>
                        {children}
                      </code>
                    );
                  },
                }}
              >
                {obj.body}
              </ReactMarkdown>
            </section>
          );
        })}
      </article>
    </div>
  );
}

export default Blog;
