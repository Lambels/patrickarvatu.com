function Project({ project }) {
  return (
    <div className="rounded overflow-hidden shadow-lg">
      <img className="w-full" src={`http://localhost:3000/image/${project.name}_PROJECT.jpg`} alt="Project Image" />
      <div className="px-6 py-4">
        <div className="font-bold text-xl mb-2">{project.name}</div>
        <p className="text-gray-700 text-base">{project.description}</p>
      </div>
      <div className="px-6 pt-4 pb-2">
        {project.topics.map((tag) => {
          return (
            <span className="inline-block bg-gray-200 rounded-full px-3 py-1 text-sm font-semibold text-gray-700 mr-2 mb-2">
              {tag}
            </span>
          );
        })}
      </div>
    </div>
  );
}

export default Project