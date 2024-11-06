import { Button, Card, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import React, { useEffect } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { usePostingInfo } from "../hooks/usePostingInfo.ts";
import { useBusinessIds } from "../hooks/useBusinessIds.ts";


export const PostingInfo = () => {
    const business_ids = useBusinessIds();
    const data = usePostingInfo({ business_ids });

    if (data.length > 0) {
        return (data.map((posting, _) => (
            <div className="content">
                <Card interactive={true} >
                    <H3>Account Information</H3>
                    <FormGroup label="Name"
                        labelFor="name" >
                        <InputGroup id="name" placeholder={posting.title} />
                    </FormGroup>
                    <FormGroup label="Email"
                        labelFor="email" >
                        <InputGroup id="email" placeholder={posting.desc} />
                    </FormGroup>
                    <Button>Save changes</Button>
                </Card>
            </div>
        )));
    }
}