import { createContext, useContext, useEffect, useState } from "react";

const AuthContext = createContext({
  user: undefined,
  isAuth: false,
  logout: () => {},
  updateUser: () => {},
});

function AuthProvider({ children }) {
  const [user, setUser] = useState({});
  const [isAuth, setIsAuth] = useState(false);

  useEffect(() => {
    updateUser()
  }, [])

  const logout = () => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/v1/oauth/user/logout`, {
      method: "DELETE",
      credentials: "include",
    }).then((response) => {
      setIsAuth(!response.ok);
    });
  };

  const updateUser = () => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/v1/oauth/user/me`, {
        method: "GET",
        credentials: "include",
    }).then((response) => {
        if (!response.ok) return
        response.json()
    }).then((data) => {
        setUser(data.user)
    })
  }

  return (
    <AuthContext.Provider
      value={{
        user: user,
        isAuth: isAuth,
        logout: logout,
        updateUser: updateUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export default AuthProvider

function useAuth() {
  return useContext(AuthContext);
}
