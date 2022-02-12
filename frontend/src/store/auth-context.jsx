import { useRouter } from "next/router";
import { createContext, useContext, useEffect, useState } from "react";

const AuthContext = createContext({
  user: undefined,
  isAuth: false,
  logout: () => {},
  redirectToProvider: () => {},
  updateUser: () => {},
});

function AuthProvider({ children }) {
  const [user, setUser] = useState({});
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
        response.json();
      })
      .then((data) => {
        setUser(data);
      });
  };

  return (
    <AuthContext.Provider
      value={{
        user: user,
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
