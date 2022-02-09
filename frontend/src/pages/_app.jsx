import "tailwindcss/tailwind.css";
import AuthProvider from "../store/auth-context";

function Application({ Component, pageProps }) {
  return (
    <AuthProvider>
      <Component {...pageProps} />
    </AuthProvider>
  )
}

export default Application
