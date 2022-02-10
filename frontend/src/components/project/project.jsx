function Project({ project }) {
  return (
    <a href={project.html_url} className="block max-w-lg bg-white rounded-lg border-2 border-gray-200 shadow-xl hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
      <img className="w-full h-60 pb-1 rounded-lg" src={`http://localhost:3000/image/${project.name}_PROJECT.jpg`} alt="Project Image"/>
      <div className="p-3">
        <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{project.name}</h5>
        <p className="font-normal text-gray-700 dark:text-gray-400">{project.description}</p>
      </div>
      <div className="p-3 pb-4">
        {project.topics.map((obj, i) => {
          return (
            <span key={i} className="inline-block p1 pl-2 pr-2 mr-2 rounded-xl bg-gray-700 border-gray-600 border text-white font-semibold">{"#"+obj}</span>
          )
        })}
      </div>
    </a>
  );
}

export default Project