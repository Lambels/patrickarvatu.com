import ProjectLayout from "../components/projects-layout";
import Project from "../components/project";
import { NextSeo } from "next-seo";

function Projects({ data }) {
  return (
    <>
      <NextSeo title="Projects" />
      <ProjectLayout>
        {data.projects.map((obj, i) => {
          return <Project project={obj} key={i}/>;
        })}
      </ProjectLayout>
    </>
  );
}

export async function getServerSideProps() {
  const response = await fetch("http://localhost:8080/v1/projects");

  const data = await response.json();
  return { props: { data } };
}

export default Projects;
