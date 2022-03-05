import { useRouter } from "next/router";
import { createContext, useContext, useEffect, useState } from "react";

const AuthContext = createContext({
  user: {},
  isAuth: false,
  isAdmin: false,
  pfpUrl: "",
  logout: () => {},
  redirectToProvider: () => {},
  updateUser: () => {},
});

function AuthProvider({ children }) {
  const [userData, setUserData] = useState({
    user: {},
    isAuth: false,
    isAdmin: false,
  });
  const router = useRouter();

  const logout = () => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/v1/oauth/user/logout`, {
      method: "DELETE",
      credentials: "include",
    }).then((response) => {
      if (response.ok)
        setUserData({
          user: {},
          isAuth: false,
          isAdmin: false,
          pfpUrl: "",
        });
    });
  };

  const redirectToProvider = (provider) => {
    router.push(`${process.env.NEXT_PUBLIC_API_URL}/v1/oauth/${provider}`);
  };

  const updateUser = async () => {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/v1/oauth/user/me`,
      {
        method: "GET",
        credentials: "include",
      }
    );

    if (!response.ok) return;
    const data = await response.json();

    setUserData({
      user: data?.user,
      isAuth: true,
      isAdmin: data?.user?.isAdmin,
      pfpUrl: data?.pfpUrl,
    })
  };

  useEffect(() => {
    updateUser();
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user: userData.user,
        isAuth: userData.isAuth,
        isAdmin: userData.isAdmin,
        pfpUrl: userData.pfpUrl,
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
