import { DefaultSeo } from "next-seo";
import "tailwindcss/tailwind.css";
import AppLayout from "../components/app-layout";
import AuthProvider from "../store/auth-context";

function Application({ Component, pageProps }) {
  return (
    <>
      <DefaultSeo
        defaultTitle="Patrick Arvatu"
        titleTemplate="%s @ Patrick Arvatu"
        description="Passionate, open-minded and outgoing backend web developer."
      />
      <AuthProvider>
        <AppLayout>
          <Component {...pageProps} />
        </AppLayout>
      </AuthProvider>
    </>
  );
}

export default Application;
