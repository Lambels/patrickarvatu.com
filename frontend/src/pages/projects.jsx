import ProjectLayout from "../components/projects-layout";
import Project from "../components/project";

function Projects({ data }) {
    return (
        <ProjectLayout>
            {data.map((obj) => {
                return <Project project={obj}/>
            })}
        </ProjectLayout>
    )
}

export async function getServerSideProps() {
    const response = await fetch("https://api.github.com/users/Lambels/repos");

    const data = await response.json();
    return { props: { data } }
}

export default Projects