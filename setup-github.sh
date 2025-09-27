#!/bin/bash

echo "🚀 Setting up GitHub repository..."

# Add the remote origin (replace 'aofiee' with actual username if different)
git remote add origin https://github.com/aofiee/testcopilot.git

# Set the default branch name to main
git branch -M main

# Push the code to GitHub
git push -u origin main

echo "✅ Repository created and code pushed to GitHub!"
echo "🌐 Repository URL: https://github.com/aofiee/testcopilot"