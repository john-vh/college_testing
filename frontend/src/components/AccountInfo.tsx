import { Button, Card, Checkbox, Colors, FormGroup, H2, H5, Icon, InputGroup } from "@blueprintjs/core";
import useAccountInfo, { useGetRole } from "../hooks/useAccountInfo.ts";
import React, { useEffect } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useIsFounder } from "../hooks/useBusinessInfo.ts";
import { Role } from "./InfoPage.tsx";
import { useLogout } from "../hooks/useLogout.ts";
import { useNavigate } from "react-router-dom";

export const AccountInfo = () => {

  const data = useAccountInfo();
  const role = useGetRole();
  const logout = useLogout();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/");
  }



  if (data != null) {
    return (
      <div className="content">
        <H2 style={{ marginBottom: "0px" }}>Account Information</H2>
        <Card interactive={true} >
          <FormGroup inline label="Name"
            labelFor="name" >
            <InputGroup id="name" defaultValue={data.name} />
          </FormGroup>
          <FormGroup inline label="Email"
            labelFor="email" >
            <InputGroup id="email" defaultValue={data.email} />
          </FormGroup>
          <FormGroup label="Roles"
            labelFor="roles" >
            <Checkbox checked={true}>User</Checkbox>
            <Checkbox checked={role === Role.Founder || role === Role.Admin}>Founder</Checkbox>
            <Checkbox checked={role === Role.Admin}>Admin</Checkbox>
          </FormGroup>
          <Button onClick={() => handleLogout()}>Logout</Button>
        </Card>
      </div >
    );
  }
}

export default AccountInfo;
