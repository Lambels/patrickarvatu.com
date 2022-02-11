import { NextSeo } from "next-seo";
import { Fade } from "react-awesome-reveal";

function Home() {
  return (
    <>
      <NextSeo title="Home" />
      <div className="mx-auto text-center text-white">
        <section className="mt-20">
          <Fade direction="up" triggerOnce cascade duration={2000}>
            <h1 className="text-6xl">
              Hello, I am{" "}
              <span className="font-extrabold text-transparent bg-clip-text bg-gradient-to-br from-blue-600 to-purple-500">
                Patrick Arvatu
              </span>
            </h1>
          </Fade>
        </section>
        <section className="mt-5">
          <Fade direction="up" triggerOnce cascade delay={1000} duration={2000}>
            <h1 className="text-2xl">
              <span className="font-extrabold text-transparent bg-clip-text bg-gradient-to-br from-blue-500 to-purple-400">
                Passionate
              </span>
              ,{" "}
              <span className="font-extrabold text-transparent bg-clip-text bg-gradient-to-br from-blue-400 to-purple-300">
                open-minded
              </span>{" "}
              and{" "}
              <span className="font-extrabold text-transparent bg-clip-text bg-gradient-to-br from-blue-200 to-purple-400">
                outgoing
              </span>{" "}
              backend web developer.
            </h1>
          </Fade>
        </section>
        <section className="mt-64 flex justify-around">
          <Fade
            direction="down"
            triggerOnce
            cascade
            delay={1000}
            duration={2000}
          >
            <div className="block max-w-md bg-white rounded-lg border-2 border-gray-200 shadow-xl hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700 transition ease-in-out delay-150 hover:-translate-y-1 hover:scale-109 duration-300">
              <div className="p-3">
                <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
                  Day to day life
                </h5>
                <p className="font-normal text-gray-700 dark:text-gray-400">
                  Im a bussy guy, I play hockey professionally and aim to keep
                  my grades up, when Im not coding im either playing sports or
                  studying. I also love reading when I have time ;)
                </p>
              </div>
            </div>
            <div className="block max-w-md bg-white rounded-lg border-2 border-gray-200 shadow-xl hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700 transition ease-in-out delay-150 hover:-translate-y-1 hover:scale-109 duration-300">
              <div className="p-3">
                <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
                  Frontend Web Developing
                </h5>
                <p className="font-normal text-gray-700 dark:text-gray-400">
                  I usually do frontend web developing to see my backend in
                  action, although im not the biggest fan of it, when I do it I
                  choose react.js with next.js .
                </p>
              </div>
            </div>
            <div className="block max-w-md bg-white rounded-lg border-2 border-gray-200 shadow-xl hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700 transition ease-in-out delay-150 hover:-translate-y-1 hover:scale-109 duration-300">
              <div className="p-3">
                <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
                  Backend Web Developing
                </h5>
                <p className="font-normal text-gray-700 dark:text-gray-400">
                  Backend Web Developing is my current main focus, I work on
                  backend using the beautiful golang to produce fast and
                  ideomatic code.
                </p>
              </div>
            </div>
          </Fade>
        </section>
      </div>
    </>
  );
}

export default Home;
