import { DefaultSeo } from "next-seo";
import "tailwindcss/tailwind.css";
import AppLayout from "../components/app-layout";
import Header from "../components/header";
import AuthProvider from "../store/auth-context";

const bgLookup = {
  0: {
    bg: "bg-[url('/waves0.svg')]",
  },
  1: {
    bg: "bg-[url('/waves1.svg')]",
  },
  2: {
    bg: "bg-[url('/waves2.svg')]",
  },
  3: {
    bg: "bg-[url('/waves3.svg')]",
  },
  4: {
    bg: "bg-[url('/waves4.svg')]",
  }
};

function Application({ Component, pageProps }) {
  const randNum = Math.floor(Math.random() * 5);

  return (
    <>
      <DefaultSeo
        defaultTitle="Patrick Arvatu"
        titleTemplate="%s @ Patrick Arvatu"
        description="Passionate, open-minded and outgoing backend web developer."
      />
      <AuthProvider>
        <AppLayout bgConf={bgLookup[randNum]}>
          <Header />
          <Component {...pageProps} />
        </AppLayout>
      </AuthProvider>
    </>
  );
}

export default Application;
