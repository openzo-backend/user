#!/bin/bash



# Create a backup branch
echo "Creating a backup branch..."
git branch backup
git push origin backup

# Step 2: Remove History
echo "Creating a new orphan branch to clean history..."
git checkout --orphan latest_branch



# Remove cached files from the index
echo "Removing cached files from index..."
git rm -r --cached .

# Step 3: Add .gitignore
echo "Adding .gitignore file..."
cat <<EOL > .gitignore
/config/*.yaml
*.properties
test.db
clean.sh
/vendor/
/.gitignore
EOL

# Stage all files
echo "Staging all files..."
git add -A

# Commit the changes
echo "Committing the clean history..."
git commit -am "Initial commit - Clean history"

# Delete the old main branch
echo "Deleting the old main branch..."
git branch -D main

# Rename the new branch to main
echo "Renaming the new branch to main..."
git branch -m main

# # Remove cached files from the index
# echo "Removing cached files from index..."
# git rm -r --cached .

# Force push the new main branch to overwrite the repository history
echo "Force pushing the new main branch..."
git push -f origin main



echo "Cleanup complete. Ensure to rotate any exposed credentials."
