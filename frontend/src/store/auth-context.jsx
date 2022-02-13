import { useRouter } from "next/router";
import { createContext, useContext, useEffect, useState } from "react";

const AuthContext = createContext({
  data: undefined,
  isAuth: false,
  logout: () => {},
  redirectToProvider: () => {},
  updateUser: () => {},
});

function AuthProvider({ children }) {
  const [data, setData] = useState({});
  const [isAuth, setIsAuth] = useState(false);
  const router = useRouter();

  useEffect(() => {
    updateUser();
  }, []);

  const logout = () => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/v1/oauth/user/logout`, {
      method: "DELETE",
      credentials: "include",
    }).then((response) => {
      if (response.ok) setIsAuth(false);
    });
  };

  const redirectToProvider = (provider) => {
    router.push(`${process.env.NEXT_PUBLIC_API_URL}/v1/oauth/${provider}`);
  };

  const updateUser = () => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/v1/oauth/user/me`, {
      method: "GET",
      credentials: "include",
    })
      .then((response) => {
        if (!response.ok) return;
        return response.json();
      })
      .then((data) => {
        setData(data);
        if (data !== undefined) setIsAuth(true);
      });
  };

  return (
    <AuthContext.Provider
      value={{
        data: data,
        isAuth: isAuth,
        logout: logout,
        redirectToProvider: redirectToProvider,
        updateUser: updateUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export default AuthProvider;

export function useAuth() {
  return useContext(AuthContext);
}
