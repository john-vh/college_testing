import { Button, Card, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import React, { useEffect, useState } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { usePostingInfo } from "../hooks/usePostingInfo.ts";
import { useBusinessIds } from "../hooks/useBusinessIds.ts";
import { AddPosting } from "./AddPosting.tsx";


export const PostingInfo = () => {
    const data = usePostingInfo();
    const business_options = new Set(data.map((post) => post.business_id));

    const [postingAdd, setPostingAdd] = useState(false);

    const addNewPosting = () => {
        setPostingAdd(true); // Change state to show the Add Business page
    };

    if (postingAdd) {
        return <AddPosting businesses={business_options} setPostingAdd={setPostingAdd} />
    }

    if (data.length > 0) {
        return (
            <div style={{ width: '100%', padding: '20px' }}>
                <Button
                    intent="primary"
                    large={true}
                    style={{ width: "100%", marginBottom: "20px" }}
                    onClick={() => addNewPosting()}
                >
                    Add Posting
                </Button>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                    {data.map((posting, _) => (
                        <Card interactive={true} >
                            <H3>Posting Information</H3>
                            <FormGroup label="Title"
                                labelFor="title" >
                                <InputGroup id="title" placeholder={posting.title} />
                            </FormGroup>
                            <FormGroup label="Description"
                                labelFor="desc" >
                                <InputGroup id="desc" placeholder={posting.desc} />
                            </FormGroup>
                            <FormGroup label="Business"
                                labelFor="business" >
                                <InputGroup id="business" placeholder={posting.business_id} />
                            </FormGroup>
                            <Button>Manage posting</Button>
                        </Card>
                    ))}
                </div>
            </div>);
    }
}